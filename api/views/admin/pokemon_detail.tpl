<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta
        name="viewport"
        content="width=device-width, initial-scale=1.0"
    >

    <title>
        {{if .ViewData.PokemonSpecies}}
            {{.ViewData.PokemonSpecies.NameJA}} |
        {{end}}
        ポケモン詳細 | pokedame
    </title>

    <link rel="stylesheet" href="/static/css/main.css">
</head>

<body>
<main class="admin-layout">
    <header class="admin-header">
        <a
            class="admin-back-link"
            href="/admin/pokemon"
        >
            ← ポケモン一覧へ戻る
        </a>

        {{if .ViewData.PokemonSpecies}}
            <h1 class="admin-header__title">
                {{.ViewData.PokemonSpecies.NameJA}}
            </h1>

            <p class="admin-header__description">
                登録されているポケモン情報を表示します。
            </p>
        {{else}}
            <h1 class="admin-header__title">
                ポケモン詳細
            </h1>
        {{end}}
    </header>

    {{if .ViewData.ErrorMessage}}
        <div class="admin-message admin-message--visible admin-message--error">
            {{.ViewData.ErrorMessage}}
        </div>
    {{end}}

    {{if .ViewData.PokemonSpecies}}
        <section class="admin-panel">
            <h2 class="admin-panel__title">
                種族情報
            </h2>

            <dl class="detail-list">
                <div class="detail-list__item">
                    <dt class="detail-list__label">
                        DB ID
                    </dt>

                    <dd class="detail-list__value">
                        {{.ViewData.PokemonSpecies.ID}}
                    </dd>
                </div>

                <div class="detail-list__item">
                    <dt class="detail-list__label">
                        全国図鑑番号
                    </dt>

                    <dd class="detail-list__value">
                        No.{{.ViewData.PokemonSpecies.NationalDexNumber}}
                    </dd>
                </div>

                <div class="detail-list__item">
                    <dt class="detail-list__label">
                        日本語名
                    </dt>

                    <dd class="detail-list__value">
                        {{.ViewData.PokemonSpecies.NameJA}}
                    </dd>
                </div>

                <div class="detail-list__item">
                    <dt class="detail-list__label">
                        英語名
                    </dt>

                    <dd class="detail-list__value">
                        {{if .ViewData.PokemonSpecies.NameEN}}
                            {{.ViewData.PokemonSpecies.NameEN}}
                        {{else}}
                            未登録
                        {{end}}
                    </dd>
                </div>

                <div class="detail-list__item">
                    <dt class="detail-list__label">
                        slug
                    </dt>

                    <dd class="detail-list__value">
                        {{.ViewData.PokemonSpecies.Slug}}
                    </dd>
                </div>

                <div class="detail-list__item">
                    <dt class="detail-list__label">
                        状態
                    </dt>

                    <dd class="detail-list__value">
                        {{if .ViewData.PokemonSpecies.IsActive}}
                            有効
                        {{else}}
                            無効
                        {{end}}
                    </dd>
                </div>
            </dl>
        </section>

        <section class="admin-panel section">
            <h2 class="admin-panel__title">
                フォーム情報
            </h2>

            {{if .ViewData.PokemonForms}}
                <div class="table-wrapper">
                    <table class="data-table">
                        <thead>
                        <tr>
                            <th>ID</th>
                            <th>フォームキー</th>
                            <th>日本語名</th>
                            <th>英語名</th>
                            <th>デフォルト</th>
                            <th>状態</th>
                        </tr>
                        </thead>

                        <tbody>
                        {{range .ViewData.PokemonForms}}
                            <tr>
                                <td>{{.ID}}</td>
                                <td>{{.FormKey}}</td>

                                <td>
                                    {{if .NameJA}}
                                        {{.NameJA}}
                                    {{else}}
                                        未登録
                                    {{end}}
                                </td>

                                <td>
                                    {{if .NameEN}}
                                        {{.NameEN}}
                                    {{else}}
                                        未登録
                                    {{end}}
                                </td>

                                <td>
                                    {{if .IsDefault}}
                                        はい
                                    {{else}}
                                        いいえ
                                    {{end}}
                                </td>

                                <td>
                                    {{if .IsActive}}
                                        有効
                                    {{else}}
                                        無効
                                    {{end}}
                                </td>
                            </tr>
                        {{end}}
                        </tbody>
                    </table>
                </div>
            {{else}}
                <p class="empty-message">
                    フォーム情報は登録されていません。
                </p>
            {{end}}
        </section>

        <section class="admin-panel section">
            <h2 class="admin-panel__title">
                素材情報
            </h2>

            {{if .ViewData.PokemonAssets}}
                <div class="asset-list">
                    {{range .ViewData.PokemonAssets}}
                        <article class="asset-card">
                            <dl class="detail-list">
                                <div class="detail-list__item">
                                    <dt class="detail-list__label">
                                        ID
                                    </dt>

                                    <dd class="detail-list__value">
                                        {{.ID}}
                                    </dd>
                                </div>

                                <div class="detail-list__item">
                                    <dt class="detail-list__label">
                                        素材種別
                                    </dt>

                                    <dd class="detail-list__value">
                                        {{.AssetType}}
                                    </dd>
                                </div>

                                <div class="detail-list__item">
                                    <dt class="detail-list__label">
                                        保存先
                                    </dt>

                                    <dd class="detail-list__value">
                                        {{.StoragePath}}
                                    </dd>
                                </div>
                            </dl>

                            {{if eq .AssetType "image"}}
                                <img
                                    class="asset-card__image"
                                    src="{{.StoragePath}}"
                                    alt="登録済みポケモン素材"
                                >
                            {{end}}
                        </article>
                    {{end}}
                </div>
            {{else}}
                <p class="empty-message">
                    素材情報は登録されていません。
                </p>
            {{end}}
        </section>
    {{else}}
        <section class="admin-panel">
            <p class="empty-message">
                ポケモン情報を表示できませんでした。
            </p>
        </section>
    {{end}}
</main>
</body>
</html>