import { useState } from 'react'

import { BattleSideForm } from './components/battle/BattleSideForm'
import { FieldSettings } from './components/battle/FieldSettings'

function App() {
  const [weather, setWeather] = useState('none')
  const [field, setField] = useState('none')

  return (
    <main className="damage-calculator">
      <header className="damage-calculator__header">
        <h1 className="damage-calculator__title">pokedame</h1>
      </header>

      <section
        className="battle-stage"
        aria-labelledby="battle-stage-title"
      >
        <h2
          id="battle-stage-title"
          className="section-heading"
        >
          バトルアニメーション
        </h2>

        <div className="battle-stage__canvas">
          <div className="battle-stage__pokemon battle-stage__pokemon--attacker">
            攻撃側ポケモン
          </div>

          <div className="battle-stage__pokemon battle-stage__pokemon--defender">
            受け側ポケモン
          </div>
        </div>
      </section>

      <section
        className="damage-result"
        aria-labelledby="damage-result-title"
      >
        <h2
          id="damage-result-title"
          className="section-heading"
        >
          ダメージ計算結果
        </h2>

        <div className="damage-result__body">
          <p className="damage-result__message">
            計算結果がここに表示されます。
          </p>
        </div>
      </section>

      <section
        className="battle-settings"
        aria-label="攻撃側と受け側の設定"
      >
        <BattleSideForm side="attacker" />
        <BattleSideForm side="defender" />
      </section>

      <FieldSettings
        weather={weather}
        field={field}
        onWeatherChange={setWeather}
        onFieldChange={setField}
      />
    </main>
  )
}

export default App