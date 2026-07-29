type FieldSettingsProps = {
    weather: string
    field: string
    onWeatherChange: (weather: string) => void
    onFieldChange: (field: string) => void
}

export function FieldSettings({
    weather,
    field,
    onWeatherChange,
    onFieldChange,
}: FieldSettingsProps) {
    return (
        <section
            className="field-settings"
            aria-labelledby="field-settings-title"
        >
            <h2
                id="field-settings-title"
                className="section-heading"
            >
                天候・フィールド
            </h2>

            <div className="field-settings__body">
                <label className="field-settings__item">
                    <span className="field-settings__label">天候</span>

                    <select
                        className="field-settings__control"
                        name="weather"
                        value={weather}
                        onChange={(event) => onWeatherChange(event.target.value)}
                    >
                        <option value="none">なし</option>
                        <option value="sunny">はれ</option>
                        <option value="rain">あめ</option>
                        <option value="sandstorm">すなあらし</option>
                        <option value="snow">ゆき</option>
                    </select>
                </label>

                <label className="field-settings__item">
                    <span className="field-settings__label">フィールド</span>

                    <select
                        className="field-settings__control"
                        name="field"
                        value={field}
                        onChange={(event) => onFieldChange(event.target.value)}
                    >
                        <option value="none">なし</option>
                        <option value="electric">エレキフィールド</option>
                        <option value="grassy">グラスフィールド</option>
                        <option value="misty">ミストフィールド</option>
                        <option value="psychic">サイコフィールド</option>
                    </select>
                </label>
            </div>
        </section>
    )
}