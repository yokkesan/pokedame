package controllers

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"api-generated/database"
	"api-generated/models"

	beego "github.com/beego/beego/v2/server/web"
)

const (
	pokemonAssetRequestTimeout = 10 * time.Second
	maxPokemonAssetFileSize    = 10 << 20
	maxPokemonAssetRequestSize = maxPokemonAssetFileSize + (1 << 20)
)

var (
	errAssetFileRequired    = errors.New("asset file is required")
	errAssetFileTooLarge    = errors.New("asset file is too large")
	errInvalidAssetMIMEType = errors.New("invalid asset MIME type")
	errInvalidAssetImage    = errors.New("invalid asset image")
	errInvalidAssetField    = errors.New("invalid asset field")
)

type PokemonAssetController struct {
	beego.Controller
}

type pokemonAssetErrorResponse struct {
	Message string `json:"message"`
}

func (c *PokemonAssetController) Create() {
	pokemonFormID, err := c.getPokemonFormID()
	if err != nil {
		c.writePokemonAssetError(
			http.StatusBadRequest,
			"ポケモンフォームIDが不正です。",
		)
		return
	}

	contentType := c.Ctx.Request.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil || mediaType != "multipart/form-data" {
		c.writePokemonAssetError(
			http.StatusUnsupportedMediaType,
			"Content-Typeにはmultipart/form-dataを指定してください。",
		)
		return
	}

	c.Ctx.Request.Body = http.MaxBytesReader(
		c.Ctx.ResponseWriter,
		c.Ctx.Request.Body,
		maxPokemonAssetRequestSize,
	)

	file, fileHeader, err := c.GetFile("file")
	if err != nil {
		c.handlePokemonAssetUploadError(err)
		return
	}
	defer file.Close()

	upload, err := c.buildPokemonAssetUpload(
		pokemonFormID,
		file,
		fileHeader,
	)
	if err != nil {
		c.handlePokemonAssetUploadError(err)
		return
	}

	storagePath, absolutePath, err := savePokemonAssetFile(
		upload.data,
		upload.extension,
	)
	if err != nil {
		c.writePokemonAssetError(
			http.StatusInternalServerError,
			"素材ファイルを保存できませんでした。",
		)
		return
	}

	db, err := database.Get()
	if err != nil {
		_ = os.Remove(absolutePath)

		c.writePokemonAssetError(
			http.StatusInternalServerError,
			"データベースへ接続できませんでした。",
		)
		return
	}

	ctx, cancel := context.WithTimeout(
		c.Ctx.Request.Context(),
		pokemonAssetRequestTimeout,
	)
	defer cancel()

	asset, err := models.CreatePokemonAsset(
		ctx,
		db,
		models.CreatePokemonAssetRequest{
			PokemonFormID:    pokemonFormID,
			AssetType:        upload.assetType,
			StoragePath:      storagePath,
			OriginalFilename: upload.originalFilename,
			MimeType:         upload.mimeType,
			FileSize:         int64(len(upload.data)),
			Width:            &upload.width,
			Height:           &upload.height,
			FrameCount:       upload.frameCount,
			FrameWidth:       upload.frameWidth,
			FrameHeight:      upload.frameHeight,
			FrameRate:        upload.frameRate,
			IsLoop:           upload.isLoop,
			ChecksumSHA256:   upload.checksumSHA256,
			IsActive:         upload.isActive,
		},
	)
	if err != nil {
		_ = os.Remove(absolutePath)
		c.handlePokemonAssetCreateError(err)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = asset
	c.ServeJSON()
}

func (c *PokemonAssetController) List() {
	pokemonFormID, err := c.getPokemonFormID()
	if err != nil {
		c.writePokemonAssetError(
			http.StatusBadRequest,
			"ポケモンフォームIDが不正です。",
		)
		return
	}

	db, err := database.Get()
	if err != nil {
		c.writePokemonAssetError(
			http.StatusInternalServerError,
			"データベースへ接続できませんでした。",
		)
		return
	}

	ctx, cancel := context.WithTimeout(
		c.Ctx.Request.Context(),
		pokemonAssetRequestTimeout,
	)
	defer cancel()

	assets, err := models.FindPokemonAssetsByFormID(
		ctx,
		db,
		pokemonFormID,
	)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrPokemonFormNotFound):
			c.writePokemonAssetError(
				http.StatusNotFound,
				"指定されたポケモンフォームが存在しません。",
			)

		case errors.Is(err, models.ErrInvalidPokemonFormID):
			c.writePokemonAssetError(
				http.StatusBadRequest,
				"ポケモンフォームIDが不正です。",
			)

		default:
			c.writePokemonAssetError(
				http.StatusInternalServerError,
				"素材一覧を取得できませんでした。",
			)
		}
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = assets
	c.ServeJSON()
}

func (c *PokemonAssetController) Delete() {
	assetID, err := c.getPokemonAssetID()
	if err != nil {
		c.writePokemonAssetError(
			http.StatusBadRequest,
			"素材IDが不正です。",
		)
		return
	}

	db, err := database.Get()
	if err != nil {
		c.writePokemonAssetError(
			http.StatusInternalServerError,
			"データベースへ接続できませんでした。",
		)
		return
	}

	ctx, cancel := context.WithTimeout(
		c.Ctx.Request.Context(),
		pokemonAssetRequestTimeout,
	)
	defer cancel()

	asset, err := models.FindPokemonAssetByID(ctx, db, assetID)
	if err != nil {
		c.handlePokemonAssetDeleteError(err)
		return
	}

	absolutePath, err := resolvePokemonAssetPath(asset.StoragePath)
	if err != nil {
		c.writePokemonAssetError(
			http.StatusInternalServerError,
			"素材の保存先が不正です。",
		)
		return
	}

	if err := models.DeletePokemonAsset(ctx, db, assetID); err != nil {
		c.handlePokemonAssetDeleteError(err)
		return
	}

	if err := os.Remove(absolutePath); err != nil &&
		!errors.Is(err, os.ErrNotExist) {
		c.writePokemonAssetError(
			http.StatusInternalServerError,
			"素材情報は削除されましたが、ファイルを削除できませんでした。",
		)
		return
	}

	c.Ctx.Output.SetStatus(http.StatusNoContent)
}

type pokemonAssetUpload struct {
	data             []byte
	assetType        string
	originalFilename string
	mimeType         string
	extension        string
	width            int
	height           int
	frameCount       int
	frameWidth       *int
	frameHeight      *int
	frameRate        *float64
	isLoop           bool
	isActive         bool
	checksumSHA256   string
}

func (c *PokemonAssetController) buildPokemonAssetUpload(
	pokemonFormID int64,
	file multipart.File,
	fileHeader *multipart.FileHeader,
) (*pokemonAssetUpload, error) {
	if pokemonFormID <= 0 {
		return nil, models.ErrInvalidPokemonFormID
	}

	if fileHeader == nil {
		return nil, errAssetFileRequired
	}

	if fileHeader.Size <= 0 {
		return nil, errAssetFileRequired
	}

	if fileHeader.Size > maxPokemonAssetFileSize {
		return nil, errAssetFileTooLarge
	}

	data, err := io.ReadAll(
		io.LimitReader(file, maxPokemonAssetFileSize+1),
	)
	if err != nil {
		return nil, fmt.Errorf("read uploaded asset: %w", err)
	}

	if len(data) == 0 {
		return nil, errAssetFileRequired
	}

	if len(data) > maxPokemonAssetFileSize {
		return nil, errAssetFileTooLarge
	}

	detectedMIMEType := http.DetectContentType(data)

	extension, ok := allowedPokemonAssetExtension(detectedMIMEType)
	if !ok {
		return nil, errInvalidAssetMIMEType
	}

	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || config.Width <= 0 || config.Height <= 0 {
		return nil, errInvalidAssetImage
	}

	assetType := strings.TrimSpace(c.GetString("asset_type"))
	if !isValidPokemonAssetType(assetType) {
		return nil, models.ErrInvalidPokemonAssetType
	}

	frameCount, err := parsePositiveIntegerField(
		c.GetString("frame_count"),
		1,
	)
	if err != nil {
		return nil, errInvalidAssetField
	}

	frameWidth, err := parseOptionalPositiveIntegerField(
		c.GetString("frame_width"),
	)
	if err != nil {
		return nil, errInvalidAssetField
	}

	frameHeight, err := parseOptionalPositiveIntegerField(
		c.GetString("frame_height"),
	)
	if err != nil {
		return nil, errInvalidAssetField
	}

	frameRate, err := parseOptionalPositiveFloatField(
		c.GetString("frame_rate"),
	)
	if err != nil {
		return nil, errInvalidAssetField
	}

	isLoop, err := parseBooleanField(c.GetString("is_loop"), false)
	if err != nil {
		return nil, errInvalidAssetField
	}

	isActive, err := parseBooleanField(c.GetString("is_active"), true)
	if err != nil {
		return nil, errInvalidAssetField
	}

	originalFilename := strings.TrimSpace(
		filepath.Base(fileHeader.Filename),
	)
	if originalFilename == "" ||
		len([]rune(originalFilename)) > 255 {
		return nil, errInvalidAssetField
	}

	checksum := calculateSHA256(data)

	return &pokemonAssetUpload{
		data:             data,
		assetType:        assetType,
		originalFilename: originalFilename,
		mimeType:         detectedMIMEType,
		extension:        extension,
		width:            config.Width,
		height:           config.Height,
		frameCount:       frameCount,
		frameWidth:       frameWidth,
		frameHeight:      frameHeight,
		frameRate:        frameRate,
		isLoop:           isLoop,
		isActive:         isActive,
		checksumSHA256:   checksum,
	}, nil
}

func savePokemonAssetFile(
	data []byte,
	extension string,
) (string, string, error) {
	randomName, err := generateSecureFilename()
	if err != nil {
		return "", "", err
	}

	storagePath := filepath.Join(
		"pokemon-assets",
		randomName+extension,
	)

	absolutePath, err := resolvePokemonAssetPath(storagePath)
	if err != nil {
		return "", "", err
	}

	if err := os.MkdirAll(
		filepath.Dir(absolutePath),
		0o750,
	); err != nil {
		return "", "", fmt.Errorf(
			"create asset storage directory: %w",
			err,
		)
	}

	output, err := os.OpenFile(
		absolutePath,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL,
		0o600,
	)
	if err != nil {
		return "", "", fmt.Errorf("create asset file: %w", err)
	}

	writeSucceeded := false
	defer func() {
		_ = output.Close()

		if !writeSucceeded {
			_ = os.Remove(absolutePath)
		}
	}()

	if _, err := output.Write(data); err != nil {
		return "", "", fmt.Errorf("write asset file: %w", err)
	}

	if err := output.Sync(); err != nil {
		return "", "", fmt.Errorf("sync asset file: %w", err)
	}

	writeSucceeded = true

	return filepath.ToSlash(storagePath), absolutePath, nil
}

func resolvePokemonAssetPath(storagePath string) (string, error) {
	storageRoot := os.Getenv("POKEMON_ASSET_STORAGE_DIR")
	if storageRoot == "" {
		storageRoot = "storage"
	}

	absoluteRoot, err := filepath.Abs(storageRoot)
	if err != nil {
		return "", fmt.Errorf("resolve asset storage root: %w", err)
	}

	cleanStoragePath := filepath.Clean(
		filepath.FromSlash(storagePath),
	)
	if cleanStoragePath == "." ||
		filepath.IsAbs(cleanStoragePath) ||
		strings.HasPrefix(cleanStoragePath, "..") {
		return "", errors.New("invalid asset storage path")
	}

	absolutePath := filepath.Join(
		absoluteRoot,
		cleanStoragePath,
	)

	relativePath, err := filepath.Rel(
		absoluteRoot,
		absolutePath,
	)
	if err != nil {
		return "", fmt.Errorf("validate asset storage path: %w", err)
	}

	if relativePath == ".." ||
		strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) {
		return "", errors.New("asset storage path escapes root")
	}

	return absolutePath, nil
}

func generateSecureFilename() (string, error) {
	randomBytes := make([]byte, 16)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("generate asset filename: %w", err)
	}

	return hex.EncodeToString(randomBytes), nil
}

func allowedPokemonAssetExtension(mimeType string) (string, bool) {
	switch mimeType {
	case "image/png":
		return ".png", true

	case "image/jpeg":
		return ".jpg", true

	case "image/gif":
		return ".gif", true

	default:
		return "", false
	}
}

func isValidPokemonAssetType(assetType string) bool {
	switch assetType {
	case "image",
		"idle",
		"enter",
		"physical_attack",
		"special_attack",
		"damage",
		"faint",
		"victory":
		return true

	default:
		return false
	}
}

func parsePositiveIntegerField(
	value string,
	defaultValue int,
) (int, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return defaultValue, nil
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil || parsedValue <= 0 {
		return 0, errInvalidAssetField
	}

	return parsedValue, nil
}

func parseOptionalPositiveIntegerField(
	value string,
) (*int, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return nil, nil
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil || parsedValue <= 0 {
		return nil, errInvalidAssetField
	}

	return &parsedValue, nil
}

func parseOptionalPositiveFloatField(
	value string,
) (*float64, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return nil, nil
	}

	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil || parsedValue <= 0 {
		return nil, errInvalidAssetField
	}

	return &parsedValue, nil
}

func parseBooleanField(
	value string,
	defaultValue bool,
) (bool, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return defaultValue, nil
	}

	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, errInvalidAssetField
	}

	return parsedValue, nil
}

func calculateSHA256(data []byte) string {
	hash := sha256Sum(data)
	return hex.EncodeToString(hash)
}

func sha256Sum(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}

func (c *PokemonAssetController) getPokemonFormID() (int64, error) {
	value := c.Ctx.Input.Param(":formId")

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, models.ErrInvalidPokemonFormID
	}

	return id, nil
}

func (c *PokemonAssetController) getPokemonAssetID() (int64, error) {
	value := c.Ctx.Input.Param(":assetId")

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, models.ErrInvalidPokemonAssetID
	}

	return id, nil
}

func (c *PokemonAssetController) handlePokemonAssetUploadError(err error) {
	var maxBytesError *http.MaxBytesError

	switch {
	case errors.As(err, &maxBytesError),
		errors.Is(err, multipart.ErrMessageTooLarge),
		errors.Is(err, errAssetFileTooLarge):
		c.writePokemonAssetError(
			http.StatusRequestEntityTooLarge,
			"素材ファイルは10MB以下にしてください。",
		)

	case errors.Is(err, errAssetFileRequired):
		c.writePokemonAssetError(
			http.StatusBadRequest,
			"素材ファイルを指定してください。",
		)

	case errors.Is(err, errInvalidAssetMIMEType):
		c.writePokemonAssetError(
			http.StatusUnsupportedMediaType,
			"PNG、JPEG、GIF形式の画像のみアップロードできます。",
		)

	case errors.Is(err, errInvalidAssetImage):
		c.writePokemonAssetError(
			http.StatusBadRequest,
			"画像ファイルの内容が不正です。",
		)

	case errors.Is(err, models.ErrInvalidPokemonAssetType):
		c.writePokemonAssetError(
			http.StatusBadRequest,
			"素材種別が不正です。",
		)

	case errors.Is(err, errInvalidAssetField):
		c.writePokemonAssetError(
			http.StatusBadRequest,
			"素材情報の入力値が不正です。",
		)

	default:
		c.writePokemonAssetError(
			http.StatusBadRequest,
			"アップロードリクエストが不正です。",
		)
	}
}

func (c *PokemonAssetController) handlePokemonAssetCreateError(err error) {
	switch {
	case errors.Is(err, models.ErrInvalidPokemonFormID):
		c.writePokemonAssetError(
			http.StatusBadRequest,
			"ポケモンフォームIDが不正です。",
		)

	case errors.Is(err, models.ErrPokemonFormNotFound):
		c.writePokemonAssetError(
			http.StatusNotFound,
			"指定されたポケモンフォームが存在しません。",
		)

	case errors.Is(err, models.ErrInvalidPokemonAssetType):
		c.writePokemonAssetError(
			http.StatusBadRequest,
			"素材種別が不正です。",
		)

	case errors.Is(err, models.ErrDuplicateAssetStoragePath):
		c.writePokemonAssetError(
			http.StatusConflict,
			"同じ保存先の素材が既に登録されています。",
		)

	default:
		c.writePokemonAssetError(
			http.StatusInternalServerError,
			"素材を登録できませんでした。",
		)
	}
}

func (c *PokemonAssetController) handlePokemonAssetDeleteError(err error) {
	switch {
	case errors.Is(err, models.ErrInvalidPokemonAssetID):
		c.writePokemonAssetError(
			http.StatusBadRequest,
			"素材IDが不正です。",
		)

	case errors.Is(err, models.ErrPokemonAssetNotFound):
		c.writePokemonAssetError(
			http.StatusNotFound,
			"指定された素材が存在しません。",
		)

	default:
		c.writePokemonAssetError(
			http.StatusInternalServerError,
			"素材を削除できませんでした。",
		)
	}
}

func (c *PokemonAssetController) writePokemonAssetError(
	status int,
	message string,
) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = pokemonAssetErrorResponse{
		Message: message,
	}
	c.ServeJSON()
}
