<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta
        name="viewport"
        content="width=device-width, initial-scale=1.0"
    >
    <title>ポケモン管理 | pokedame</title>

    <style>
        :root {
            color-scheme: dark;
            font-family:
                Inter,
                "Hiragino Kaku Gothic ProN",
                "Yu Gothic",
                sans-serif;
            background: #111827;
            color: #f9fafb;
        }

        * {
            box-sizing: border-box;
        }

        body {
            margin: 0;
            background: #111827;
        }

        button,
        input,
        select {
            font: inherit;
        }

        button {
            cursor: pointer;
        }

        .admin-layout {
            width: min(1200px, calc(100% - 32px));
            margin: 0 auto;
            padding: 32px 0 64px;
        }

        .admin-header {
            margin-bottom: 24px;
        }

        .admin-header__title {
            margin: 0;
            font-size: 28px;
        }

        .admin-header__description {
            margin: 8px 0 0;
            color: #9ca3af;
        }

        .admin-message {
            display: none;
            margin-bottom: 20px;
            padding: 12px 16px;
            border-radius: 8px;
        }

        .admin-message--visible {
            display: block;
        }

        .admin-message--success {
            background: #064e3b;
            color: #d1fae5;
        }

        .admin-message--error {
            background: #7f1d1d;
            color: #fee2e2;
        }

        .admin-grid {
            display: grid;
            grid-template-columns: minmax(260px, 0.8fr) minmax(0, 1.7fr);
            gap: 20px;
        }

        .admin-panel {
            padding: 20px;
            border: 1px solid #374151;
            border-radius: 12px;
            background: #1f2937;
        }

        .admin-panel__title {
            margin: 0 0 16px;
            font-size: 20px;
        }

        .species-list {
            display: grid;
            gap: 8px;
        }

        .species-button {
            width: 100%;
            padding: 12px;
            border: 1px solid #4b5563;
            border-radius: 8px;
            background: #111827;
            color: #f9fafb;
            text-align: left;
        }

        .species-button:hover {
            border-color: #60a5fa;
        }

        .species-button--selected {
            border-color: #3b82f6;
            background: #1e3a8a;
        }

        .species-button__name {
            display: block;
            font-weight: 700;
        }

        .species-button__meta {
            display: block;
            margin-top: 4px;
            color: #9ca3af;
            font-size: 13px;
        }

        .section {
            margin-top: 28px;
        }

        .section:first-child {
            margin-top: 0;
        }

        .section__header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 16px;
            margin-bottom: 12px;
        }

        .section__title {
            margin: 0;
            font-size: 18px;
        }

        .form-grid {
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));
            gap: 12px;
        }

        .form-field {
            display: grid;
            gap: 6px;
        }

        .form-field--full {
            grid-column: 1 / -1;
        }

        .form-field__label {
            color: #d1d5db;
            font-size: 14px;
        }

        .form-field__input {
            width: 100%;
            padding: 10px 12px;
            border: 1px solid #4b5563;
            border-radius: 8px;
            background: #111827;
            color: #f9fafb;
        }

        .form-field__input:focus {
            border-color: #60a5fa;
            outline: none;
        }

        .checkbox-field {
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .button {
            padding: 10px 16px;
            border: 0;
            border-radius: 8px;
            font-weight: 700;
        }

        .button--primary {
            background: #2563eb;
            color: #ffffff;
        }

        .button--danger {
            background: #b91c1c;
            color: #ffffff;
        }

        .button:disabled {
            cursor: not-allowed;
            opacity: 0.5;
        }

        .table-wrapper {
            overflow-x: auto;
        }

        .data-table {
            width: 100%;
            border-collapse: collapse;
        }

        .data-table th,
        .data-table td {
            padding: 10px;
            border-bottom: 1px solid #374151;
            text-align: left;
            white-space: nowrap;
        }

        .data-table th {
            color: #9ca3af;
            font-size: 13px;
        }

        .empty-message {
            margin: 0;
            padding: 20px 0;
            color: #9ca3af;
            text-align: center;
        }

        .preview-image {
            width: 72px;
            height: 72px;
            border-radius: 8px;
            background: #111827;
            object-fit: contain;
        }

        @media (max-width: 800px) {
            .admin-grid,
            .form-grid {
                grid-template-columns: 1fr;
            }

            .form-field--full {
                grid-column: auto;
            }
        }
    </style>
</head>

<body>
<main class="admin-layout">
    <header class="admin-header">
        <h1 class="admin-header__title">ポケモン管理</h1>
        <p class="admin-header__description">
            種族ごとのフォームと画像・アニメーション素材を管理します。
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
            <h2 class="admin-panel__title">ポケモン種族</h2>

            <div class="species-list">
                {{range .ViewData.PokemonSpecies}}
                    <button
                        type="button"
                        class="species-button"
                        data-species-id="{{.ID}}"
                        data-species-name="{{.NameJA}}"
                    >
                        <span class="species-button__name">
                            {{.NameJA}}
                        </span>

                        <span class="species-button__meta">
                            No.{{.NationalDexNumber}} / {{.Slug}}
                        </span>
                    </button>
                {{else}}
                    <p class="empty-message">
                        ポケモン種族が登録されていません。
                    </p>
                {{end}}
            </div>
        </section>

        <section class="admin-panel">
            <div class="section">
                <div class="section__header">
                    <h2
                        id="selected-species-title"
                        class="section__title"
                    >
                        フォーム管理
                    </h2>
                </div>

                <form id="pokemon-form-create-form">
                    <div class="form-grid">
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
                                disabled
                            >
                                フォームを登録
                            </button>
                        </div>
                    </div>
                </form>

                <div class="table-wrapper">
                    <table class="data-table">
                        <thead>
                        <tr>
                            <th>ID</th>
                            <th>フォームキー</th>
                            <th>日本語名</th>
                            <th>デフォルト</th>
                            <th>状態</th>
                            <th>素材</th>
                        </tr>
                        </thead>

                        <tbody id="pokemon-form-list">
                        <tr>
                            <td colspan="6">
                                <p class="empty-message">
                                    左側からポケモン種族を選択してください。
                                </p>
                            </td>
                        </tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <div class="section">
                <div class="section__header">
                    <h2
                        id="asset-section-title"
                        class="section__title"
                    >
                        素材管理
                    </h2>
                </div>

                <form id="pokemon-asset-upload-form">
                    <div class="form-grid">
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

                <div class="table-wrapper">
                    <table class="data-table">
                        <thead>
                        <tr>
                            <th>ID</th>
                            <th>種別</th>
                            <th>ファイル名</th>
                            <th>サイズ</th>
                            <th>画像サイズ</th>
                            <th>操作</th>
                        </tr>
                        </thead>

                        <tbody id="pokemon-asset-list">
                        <tr>
                            <td colspan="6">
                                <p class="empty-message">
                                    フォームを選択してください。
                                </p>
                            </td>
                        </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </section>
    </div>
</main>

<script>
    (() => {
        "use strict";

        const state = {
            speciesId: null,
            formId: null,
        };

        const messageElement =
            document.getElementById("admin-message");

        const formCreateForm =
            document.getElementById("pokemon-form-create-form");

        const assetUploadForm =
            document.getElementById("pokemon-asset-upload-form");

        const formListElement =
            document.getElementById("pokemon-form-list");

        const assetListElement =
            document.getElementById("pokemon-asset-list");

        const createFormButton =
            document.getElementById("create-form-button");

        const uploadAssetButton =
            document.getElementById("upload-asset-button");

        const selectedSpeciesTitle =
            document.getElementById("selected-species-title");

        const assetSectionTitle =
            document.getElementById("asset-section-title");

        function escapeHTML(value) {
            const element = document.createElement("div");
            element.textContent = String(value ?? "");
            return element.innerHTML;
        }

        function showMessage(message, type) {
            messageElement.textContent = message;
            messageElement.className =
                `admin-message admin-message--visible admin-message--${type}`;
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
                throw new Error("サーバーから不正な応答が返されました。");
            }
        }

        async function loadForms() {
            if (!state.speciesId) {
                return;
            }

            const response = await fetch(
                `/api/admin/pokemon-species/${state.speciesId}/forms`,
                {
                    method: "GET",
                    headers: {
                        "Accept": "application/json",
                    },
                },
            );

            const data = await readJSON(response);

            if (!response.ok) {
                throw new Error(
                    data?.message ?? "フォーム一覧を取得できませんでした。",
                );
            }

            renderForms(data);
        }

        function renderForms(forms) {
            state.formId = null;
            uploadAssetButton.disabled = true;
            assetSectionTitle.textContent = "素材管理";

            assetListElement.innerHTML = `
                <tr>
                    <td colspan="6">
                        <p class="empty-message">
                            フォームを選択してください。
                        </p>
                    </td>
                </tr>
            `;

            if (!Array.isArray(forms) || forms.length === 0) {
                formListElement.innerHTML = `
                    <tr>
                        <td colspan="6">
                            <p class="empty-message">
                                フォームが登録されていません。
                            </p>
                        </td>
                    </tr>
                `;
                return;
            }

            formListElement.innerHTML = forms.map((form) => `
                <tr>
                    <td>${form.id}</td>
                    <td>${escapeHTML(form.form_key)}</td>
                    <td>${escapeHTML(form.name_ja ?? "-")}</td>
                    <td>${form.is_default ? "はい" : "いいえ"}</td>
                    <td>${form.is_active ? "有効" : "無効"}</td>
                    <td>
                        <button
                            type="button"
                            class="button button--primary"
                            data-action="select-form"
                            data-form-id="${form.id}"
                            data-form-name="${escapeHTML(
                                form.name_ja ?? form.form_key,
                            )}"
                        >
                            素材を管理
                        </button>
                    </td>
                </tr>
            `).join("");
        }

        async function loadAssets() {
            if (!state.formId) {
                return;
            }

            const response = await fetch(
                `/api/admin/pokemon-forms/${state.formId}/assets`,
                {
                    method: "GET",
                    headers: {
                        "Accept": "application/json",
                    },
                },
            );

            const data = await readJSON(response);

            if (!response.ok) {
                throw new Error(
                    data?.message ?? "素材一覧を取得できませんでした。",
                );
            }

            renderAssets(data);
        }

        function renderAssets(assets) {
            if (!Array.isArray(assets) || assets.length === 0) {
                assetListElement.innerHTML = `
                    <tr>
                        <td colspan="6">
                            <p class="empty-message">
                                素材が登録されていません。
                            </p>
                        </td>
                    </tr>
                `;
                return;
            }

            assetListElement.innerHTML = assets.map((asset) => `
                <tr>
                    <td>${asset.id}</td>
                    <td>${escapeHTML(asset.asset_type)}</td>
                    <td>${escapeHTML(asset.original_filename)}</td>
                    <td>${formatFileSize(asset.file_size)}</td>
                    <td>
                        ${asset.width ?? "-"} × ${asset.height ?? "-"}
                    </td>
                    <td>
                        <button
                            type="button"
                            class="button button--danger"
                            data-action="delete-asset"
                            data-asset-id="${asset.id}"
                        >
                            削除
                        </button>
                    </td>
                </tr>
            `).join("");
        }

        function formatFileSize(size) {
            const bytes = Number(size);

            if (!Number.isFinite(bytes) || bytes < 0) {
                return "-";
            }

            if (bytes < 1024) {
                return `${bytes} B`;
            }

            if (bytes < 1024 * 1024) {
                return `${(bytes / 1024).toFixed(1)} KB`;
            }

            return `${(bytes / 1024 / 1024).toFixed(1)} MB`;
        }

        document
            .querySelectorAll(".species-button")
            .forEach((button) => {
                button.addEventListener("click", async () => {
                    clearMessage();

                    document
                        .querySelectorAll(".species-button")
                        .forEach((item) => {
                            item.classList.remove(
                                "species-button--selected",
                            );
                        });

                    button.classList.add(
                        "species-button--selected",
                    );

                    state.speciesId = Number(
                        button.dataset.speciesId,
                    );

                    state.formId = null;
                    createFormButton.disabled = false;
                    uploadAssetButton.disabled = true;

                    selectedSpeciesTitle.textContent =
                        `${button.dataset.speciesName}のフォーム管理`;

                    try {
                        await loadForms();
                    } catch (error) {
                        showMessage(error.message, "error");
                    }
                });
            });

        formCreateForm.addEventListener("submit", async (event) => {
            event.preventDefault();
            clearMessage();

            if (!state.speciesId) {
                showMessage(
                    "ポケモン種族を選択してください。",
                    "error",
                );
                return;
            }

            const formData = new FormData(formCreateForm);

            const payload = {
                form_key: String(
                    formData.get("form_key") ?? "",
                ),
                name_ja: String(
                    formData.get("name_ja") ?? "",
                ),
                name_en: String(
                    formData.get("name_en") ?? "",
                ),
                is_default: formData.has("is_default"),
                is_active: formData.has("is_active"),
            };

            createFormButton.disabled = true;

            try {
                const response = await fetch(
                    `/api/admin/pokemon-species/${state.speciesId}/forms`,
                    {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                            "Accept": "application/json",
                        },
                        body: JSON.stringify(payload),
                    },
                );

                const data = await readJSON(response);

                if (!response.ok) {
                    throw new Error(
                        data?.message ?? "フォームを登録できませんでした。",
                    );
                }

                formCreateForm.reset();
                formCreateForm.elements.is_active.checked = true;

                await loadForms();

                showMessage(
                    "フォームを登録しました。",
                    "success",
                );
            } catch (error) {
                showMessage(error.message, "error");
            } finally {
                createFormButton.disabled = false;
            }
        });

        formListElement.addEventListener("click", async (event) => {
            const button = event.target.closest(
                '[data-action="select-form"]',
            );

            if (!button) {
                return;
            }

            clearMessage();

            state.formId = Number(button.dataset.formId);
            uploadAssetButton.disabled = false;

            assetSectionTitle.textContent =
                `${button.dataset.formName}の素材管理`;

            try {
                await loadAssets();
            } catch (error) {
                showMessage(error.message, "error");
            }
        });

        assetUploadForm.addEventListener(
            "submit",
            async (event) => {
                event.preventDefault();
                clearMessage();

                if (!state.formId) {
                    showMessage(
                        "フォームを選択してください。",
                        "error",
                    );
                    return;
                }

                const formData = new FormData(assetUploadForm);

                formData.set(
                    "is_loop",
                    formData.has("is_loop") ? "true" : "false",
                );

                formData.set(
                    "is_active",
                    formData.has("is_active") ? "true" : "false",
                );

                uploadAssetButton.disabled = true;

                try {
                    const response = await fetch(
                        `/api/admin/pokemon-forms/${state.formId}/assets`,
                        {
                            method: "POST",
                            headers: {
                                "Accept": "application/json",
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
                    assetUploadForm.elements.frame_count.value = "1";
                    assetUploadForm.elements.is_active.checked = true;

                    await loadAssets();

                    showMessage(
                        "素材をアップロードしました。",
                        "success",
                    );
                } catch (error) {
                    showMessage(error.message, "error");
                } finally {
                    uploadAssetButton.disabled = false;
                }
            },
        );

        assetListElement.addEventListener("click", async (event) => {
            const button = event.target.closest(
                '[data-action="delete-asset"]',
            );

            if (!button) {
                return;
            }

            const confirmed = window.confirm(
                "この素材を削除しますか？",
            );

            if (!confirmed) {
                return;
            }

            clearMessage();
            button.disabled = true;

            try {
                const response = await fetch(
                    `/api/admin/pokemon-assets/${button.dataset.assetId}`,
                    {
                        method: "DELETE",
                        headers: {
                            "Accept": "application/json",
                        },
                    },
                );

                if (!response.ok) {
                    const data = await readJSON(response);

                    throw new Error(
                        data?.message ?? "素材を削除できませんでした。",
                    );
                }

                await loadAssets();

                showMessage(
                    "素材を削除しました。",
                    "success",
                );
            } catch (error) {
                showMessage(error.message, "error");
                button.disabled = false;
            }
        });
    })();
</script>
</body>
</html>