<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta
        name="viewport"
        content="width=device-width, initial-scale=1.0"
    >

    <title>ポケモン管理 | pokedame</title>

    <link rel="stylesheet" href="/static/css/main.css">
</head>

<body>
<main class="admin-layout">
    <header class="admin-header">
        <h1 class="admin-header__title">
            ポケモン管理
        </h1>

        <p class="admin-header__description">
            ポケモン種族、フォーム、素材を登録します。
        </p>
    </header>

    {{if .ViewData.ErrorMessage}}
        <div class="admin-message admin-message--visible admin-message--error">
            {{.ViewData.ErrorMessage}}
        </div>
    {{end}}

    <div
        id="admin-message"
        class="admin-message"
        role="status"
    ></div>

    <div class="admin-grid">
        <section class="admin-panel">
            <h2 class="admin-panel__title">
                ポケモン種族
            </h2>

            <div class="species-list">
                {{range .ViewData.PokemonSpecies}}
                    <a
                        class="species-button"
                        href="/admin/pokemon/{{.ID}}"
                    >
                        <span class="species-button__name">
                            {{.NameJA}}
                        </span>

                        <span class="species-button__meta">
                            No.{{.NationalDexNumber}} / {{.Slug}}
                        </span>
                    </a>
                {{else}}
                    <p class="empty-message">
                        ポケモン種族が登録されていません。
                    </p>
                {{end}}
            </div>
        </section>

        <section class="admin-panel">
            <h2 class="admin-panel__title">
                ポケモン種族を登録
            </h2>

            <form id="pokemon-species-create-form">
                <div class="form-grid">
                    <label class="form-field">
                        <span class="form-field__label">
                            全国図鑑番号
                        </span>

                        <input
                            class="form-field__input"
                            type="number"
                            name="national_dex_number"
                            min="1"
                            required
                        >
                    </label>

                    <label class="form-field">
                        <span class="form-field__label">
                            日本語名
                        </span>

                        <input
                            class="form-field__input"
                            type="text"
                            name="name_ja"
                            maxlength="100"
                            required
                        >
                    </label>

                    <label class="form-field">
                        <span class="form-field__label">
                            英語名
                        </span>

                        <input
                            class="form-field__input"
                            type="text"
                            name="name_en"
                            maxlength="100"
                        >
                    </label>

                    <label class="form-field">
                        <span class="form-field__label">
                            slug
                        </span>

                        <input
                            class="form-field__input"
                            type="text"
                            name="slug"
                            maxlength="100"
                            placeholder="pikachu"
                            required
                        >
                    </label>

                    <label class="checkbox-field">
                        <input
                            type="checkbox"
                            name="is_active"
                            checked
                        >
                        有効
                    </label>

                    <div class="form-field form-field--full">
                        <button
                            id="create-species-button"
                            class="button button--primary"
                            type="submit"
                        >
                            ポケモン種族を登録
                        </button>
                    </div>
                </div>
            </form>
        </section>
    </div>

    <section class="admin-panel section">
        <h2 class="admin-panel__title">
            フォームを登録
        </h2>

        <form id="pokemon-form-create-form">
            <div class="form-grid">
                <label class="form-field form-field--full">
                    <span class="form-field__label">
                        ポケモン種族
                    </span>

                    <select
                        id="form-species-id"
                        class="form-field__input"
                        name="species_id"
                        required
                    >
                        <option value="">
                            選択してください
                        </option>

                        {{range .ViewData.PokemonSpecies}}
                            <option value="{{.ID}}">
                                No.{{.NationalDexNumber}} {{.NameJA}}
                            </option>
                        {{end}}
                    </select>
                </label>

                <label class="form-field">
                    <span class="form-field__label">
                        フォームキー
                    </span>

                    <input
                        class="form-field__input"
                        type="text"
                        name="form_key"
                        maxlength="100"
                        placeholder="normal"
                        required
                    >
                </label>

                <label class="form-field">
                    <span class="form-field__label">
                        日本語名
                    </span>

                    <input
                        class="form-field__input"
                        type="text"
                        name="name_ja"
                        maxlength="100"
                        placeholder="通常"
                    >
                </label>

                <label class="form-field">
                    <span class="form-field__label">
                        英語名
                    </span>

                    <input
                        class="form-field__input"
                        type="text"
                        name="name_en"
                        maxlength="100"
                        placeholder="Normal"
                    >
                </label>

                <div class="form-field">
                    <span class="form-field__label">
                        設定
                    </span>

                    <label class="checkbox-field">
                        <input
                            type="checkbox"
                            name="is_default"
                        >
                        デフォルトフォーム
                    </label>

                    <label class="checkbox-field">
                        <input
                            type="checkbox"
                            name="is_active"
                            checked
                        >
                        有効
                    </label>
                </div>

                <div class="form-field form-field--full">
                    <button
                        id="create-form-button"
                        class="button button--primary"
                        type="submit"
                    >
                        フォームを登録
                    </button>
                </div>
            </div>
        </form>
    </section>

    <section class="admin-panel section">
        <h2 class="admin-panel__title">
            素材を登録
        </h2>

        <form id="pokemon-asset-upload-form">
            <div class="form-grid">
                <label class="form-field">
                    <span class="form-field__label">
                        ポケモン種族
                    </span>

                    <select
                        id="asset-species-id"
                        class="form-field__input"
                        required
                    >
                        <option value="">
                            選択してください
                        </option>

                        {{range .ViewData.PokemonSpecies}}
                            <option value="{{.ID}}">
                                No.{{.NationalDexNumber}} {{.NameJA}}
                            </option>
                        {{end}}
                    </select>
                </label>

                <label class="form-field">
                    <span class="form-field__label">
                        フォーム
                    </span>

                    <select
                        id="asset-form-id"
                        class="form-field__input"
                        name="form_id"
                        required
                        disabled
                    >
                        <option value="">
                            先にポケモン種族を選択してください
                        </option>
                    </select>
                </label>

                <label class="form-field form-field--full">
                    <span class="form-field__label">
                        画像ファイル
                    </span>

                    <input
                        class="form-field__input"
                        type="file"
                        name="file"
                        accept="image/png,image/jpeg,image/gif"
                        required
                    >
                </label>

                <label class="form-field">
                    <span class="form-field__label">
                        素材種別
                    </span>

                    <select
                        class="form-field__input"
                        name="asset_type"
                        required
                    >
                        <option value="image">image</option>
                        <option value="idle">idle</option>
                        <option value="enter">enter</option>
                        <option value="physical_attack">
                            physical_attack
                        </option>
                        <option value="special_attack">
                            special_attack
                        </option>
                        <option value="damage">damage</option>
                        <option value="faint">faint</option>
                        <option value="victory">victory</option>
                    </select>
                </label>

                <label class="form-field">
                    <span class="form-field__label">
                        フレーム数
                    </span>

                    <input
                        class="form-field__input"
                        type="number"
                        name="frame_count"
                        value="1"
                        min="1"
                        required
                    >
                </label>

                <label class="checkbox-field">
                    <input
                        type="checkbox"
                        name="is_loop"
                    >
                    ループ再生
                </label>

                <label class="checkbox-field">
                    <input
                        type="checkbox"
                        name="is_active"
                        checked
                    >
                    有効
                </label>

                <div class="form-field form-field--full">
                    <button
                        id="upload-asset-button"
                        class="button button--primary"
                        type="submit"
                        disabled
                    >
                        素材をアップロード
                    </button>
                </div>
            </div>
        </form>
    </section>
</main>

<script>
    (() => {
        "use strict";

        const messageElement =
            document.getElementById("admin-message");

        const speciesCreateForm =
            document.getElementById(
                "pokemon-species-create-form",
            );

        const pokemonFormCreateForm =
            document.getElementById(
                "pokemon-form-create-form",
            );

        const assetUploadForm =
            document.getElementById(
                "pokemon-asset-upload-form",
            );

        const createSpeciesButton =
            document.getElementById(
                "create-species-button",
            );

        const createFormButton =
            document.getElementById(
                "create-form-button",
            );

        const uploadAssetButton =
            document.getElementById(
                "upload-asset-button",
            );

        const assetSpeciesSelect =
            document.getElementById(
                "asset-species-id",
            );

        const assetFormSelect =
            document.getElementById(
                "asset-form-id",
            );

        function showMessage(message, type) {
            messageElement.textContent = message;
            messageElement.className =
                `admin-message admin-message--visible admin-message--${type}`;

            messageElement.scrollIntoView({
                behavior: "smooth",
                block: "center",
            });
        }

        function clearMessage() {
            messageElement.textContent = "";
            messageElement.className = "admin-message";
        }

        async function readJSON(response) {
            const text = await response.text();

            if (text === "") {
                return null;
            }

            try {
                return JSON.parse(text);
            } catch {
                throw new Error(
                    "サーバーから不正な応答が返されました。",
                );
            }
        }

        speciesCreateForm.addEventListener(
            "submit",
            async (event) => {
                event.preventDefault();
                clearMessage();

                const formData =
                    new FormData(speciesCreateForm);

                const payload = {
                    national_dex_number: Number(
                        formData.get(
                            "national_dex_number",
                        ),
                    ),
                    name_ja: String(
                        formData.get("name_ja") ?? "",
                    ).trim(),
                    name_en: String(
                        formData.get("name_en") ?? "",
                    ).trim(),
                    slug: String(
                        formData.get("slug") ?? "",
                    ).trim(),
                    is_active:
                        formData.has("is_active"),
                };

                createSpeciesButton.disabled = true;

                try {
                    const response = await fetch(
                        "/api/admin/pokemon-species",
                        {
                            method: "POST",
                            headers: {
                                "Content-Type":
                                    "application/json",
                                "Accept":
                                    "application/json",
                            },
                            body: JSON.stringify(payload),
                        },
                    );

                    const data = await readJSON(response);

                    if (!response.ok) {
                        throw new Error(
                            data?.message ??
                            "ポケモン種族を登録できませんでした。",
                        );
                    }

                    window.location.reload();
                } catch (error) {
                    showMessage(
                        error.message,
                        "error",
                    );
                } finally {
                    createSpeciesButton.disabled = false;
                }
            },
        );

        pokemonFormCreateForm.addEventListener(
            "submit",
            async (event) => {
                event.preventDefault();
                clearMessage();

                const formData =
                    new FormData(pokemonFormCreateForm);

                const speciesId =
                    Number(formData.get("species_id"));

                if (!speciesId) {
                    showMessage(
                        "ポケモン種族を選択してください。",
                        "error",
                    );
                    return;
                }

                const payload = {
                    form_key: String(
                        formData.get("form_key") ?? "",
                    ).trim(),
                    name_ja: String(
                        formData.get("name_ja") ?? "",
                    ).trim(),
                    name_en: String(
                        formData.get("name_en") ?? "",
                    ).trim(),
                    is_default:
                        formData.has("is_default"),
                    is_active:
                        formData.has("is_active"),
                };

                createFormButton.disabled = true;

                try {
                    const response = await fetch(
                        `/api/admin/pokemon-species/${speciesId}/forms`,
                        {
                            method: "POST",
                            headers: {
                                "Content-Type":
                                    "application/json",
                                "Accept":
                                    "application/json",
                            },
                            body: JSON.stringify(payload),
                        },
                    );

                    const data = await readJSON(response);

                    if (!response.ok) {
                        throw new Error(
                            data?.message ??
                            "フォームを登録できませんでした。",
                        );
                    }

                    pokemonFormCreateForm.reset();

                    pokemonFormCreateForm.elements
                        .is_active.checked = true;

                    showMessage(
                        "フォームを登録しました。",
                        "success",
                    );
                } catch (error) {
                    showMessage(
                        error.message,
                        "error",
                    );
                } finally {
                    createFormButton.disabled = false;
                }
            },
        );

        assetSpeciesSelect.addEventListener(
            "change",
            async () => {
                clearMessage();

                const speciesId =
                    Number(assetSpeciesSelect.value);

                assetFormSelect.disabled = true;
                uploadAssetButton.disabled = true;

                assetFormSelect.innerHTML = `
                    <option value="">
                        読み込み中です
                    </option>
                `;

                if (!speciesId) {
                    assetFormSelect.innerHTML = `
                        <option value="">
                            先にポケモン種族を選択してください
                        </option>
                    `;
                    return;
                }

                try {
                    const response = await fetch(
                        `/api/admin/pokemon-species/${speciesId}/forms`,
                        {
                            method: "GET",
                            headers: {
                                "Accept":
                                    "application/json",
                            },
                        },
                    );

                    const forms = await readJSON(response);

                    if (!response.ok) {
                        throw new Error(
                            forms?.message ??
                            "フォーム一覧を取得できませんでした。",
                        );
                    }

                    if (
                        !Array.isArray(forms) ||
                        forms.length === 0
                    ) {
                        assetFormSelect.innerHTML = `
                            <option value="">
                                フォームが登録されていません
                            </option>
                        `;
                        return;
                    }

                    assetFormSelect.innerHTML = `
                        <option value="">
                            選択してください
                        </option>
                        ${forms.map((form) => `
                            <option value="${form.id}">
                                ${form.name_ja || form.form_key}
                            </option>
                        `).join("")}
                    `;

                    assetFormSelect.disabled = false;
                } catch (error) {
                    assetFormSelect.innerHTML = `
                        <option value="">
                            取得に失敗しました
                        </option>
                    `;

                    showMessage(
                        error.message,
                        "error",
                    );
                }
            },
        );

        assetFormSelect.addEventListener(
            "change",
            () => {
                uploadAssetButton.disabled =
                    !assetFormSelect.value;
            },
        );

        assetUploadForm.addEventListener(
            "submit",
            async (event) => {
                event.preventDefault();
                clearMessage();

                const formId =
                    Number(assetFormSelect.value);

                if (!formId) {
                    showMessage(
                        "フォームを選択してください。",
                        "error",
                    );
                    return;
                }

                const formData =
                    new FormData(assetUploadForm);

                formData.delete("form_id");

                formData.set(
                    "is_loop",
                    formData.has("is_loop")
                        ? "true"
                        : "false",
                );

                formData.set(
                    "is_active",
                    formData.has("is_active")
                        ? "true"
                        : "false",
                );

                uploadAssetButton.disabled = true;

                try {
                    const response = await fetch(
                        `/api/admin/pokemon-forms/${formId}/assets`,
                        {
                            method: "POST",
                            headers: {
                                "Accept":
                                    "application/json",
                            },
                            body: formData,
                        },
                    );

                    const data = await readJSON(response);

                    if (!response.ok) {
                        throw new Error(
                            data?.message ??
                            "素材をアップロードできませんでした。",
                        );
                    }

                    assetUploadForm.reset();

                    assetUploadForm.elements
                        .frame_count.value = "1";

                    assetUploadForm.elements
                        .is_active.checked = true;

                    assetSpeciesSelect.value = "";
                    assetFormSelect.disabled = true;

                    assetFormSelect.innerHTML = `
                        <option value="">
                            先にポケモン種族を選択してください
                        </option>
                    `;

                    showMessage(
                        "素材を登録しました。",
                        "success",
                    );
                } catch (error) {
                    showMessage(
                        error.message,
                        "error",
                    );

                    uploadAssetButton.disabled = false;
                }
            },
        );
    })();
</script>
</body>
</html>