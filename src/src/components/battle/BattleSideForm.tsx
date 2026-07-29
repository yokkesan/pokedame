type BattleSide = 'attacker' | 'defender'

type BattleSideFormProps = {
    side: BattleSide
}

const sideLabels: Record<BattleSide, string> = {
    attacker: '攻撃側',
    defender: '受け側',
}

export function BattleSideForm({ side }: BattleSideFormProps) {
    const sideLabel = sideLabels[side]

    return (
        <section className={`battle-panel battle-panel--${side}`}>
            <h2 className="battle-panel__heading">{sideLabel}</h2>

            <div className="battle-panel__body">
                <div className="battle-side-form">
                    <label className="battle-side-form__field">
                        <span className="battle-side-form__label">ポケモン</span>

                        <select
                            className="battle-side-form__control"
                            name={`${side}-pokemon`}
                            defaultValue=""
                        >
                            <option value="">選択してください</option>
                        </select>
                    </label>

                    <label className="battle-side-form__field">
                        <span className="battle-side-form__label">レベル</span>

                        <input
                            className="battle-side-form__control"
                            type="number"
                            name={`${side}-level`}
                            min="1"
                            max="100"
                            defaultValue="50"
                        />
                    </label>

                    <label className="battle-side-form__field">
                        <span className="battle-side-form__label">特性</span>

                        <select
                            className="battle-side-form__control"
                            name={`${side}-ability`}
                            defaultValue=""
                        >
                            <option value="">選択してください</option>
                        </select>
                    </label>

                    <label className="battle-side-form__field">
                        <span className="battle-side-form__label">持ち物</span>

                        <select
                            className="battle-side-form__control"
                            name={`${side}-item`}
                            defaultValue=""
                        >
                            <option value="">なし</option>
                        </select>
                    </label>

                    {side === 'attacker' && (
                        <label className="battle-side-form__field">
                            <span className="battle-side-form__label">技</span>

                            <select
                                className="battle-side-form__control"
                                name="attacker-move"
                                defaultValue=""
                            >
                                <option value="">選択してください</option>
                            </select>
                        </label>
                    )}

                    <label className="battle-side-form__field">
                        <span className="battle-side-form__label">性格補正</span>

                        <select
                            className="battle-side-form__control"
                            name={`${side}-nature`}
                            defaultValue="none"
                        >
                            <option value="none">補正なし</option>
                            <option value="up">上昇補正</option>
                            <option value="down">下降補正</option>
                        </select>
                    </label>

                    <div className="battle-side-form__stats">
                        <label className="battle-side-form__field">
                            <span className="battle-side-form__label">個体値</span>

                            <input
                                className="battle-side-form__control"
                                type="number"
                                name={`${side}-iv`}
                                min="0"
                                max="31"
                                defaultValue="31"
                            />
                        </label>

                        <label className="battle-side-form__field">
                            <span className="battle-side-form__label">努力値</span>

                            <input
                                className="battle-side-form__control"
                                type="number"
                                name={`${side}-ev`}
                                min="0"
                                max="252"
                                step="4"
                                defaultValue="0"
                            />
                        </label>
                    </div>
                </div>
            </div>
        </section>
    )
}