

# Описание API виджета для бэка

Пример данных, который должен возвращать бэк:

```ts
export const mockRecapPayload: RecapPayload = {
  schemaVersion: 1, // Версия схему, в нашем случае всегда будет 1 
  meta: { // Вся дополнительная инфа по рекапу (кто / какой год / локаль)
    vertical: 'marketplace', // Вертикаль продукта (Авито Авто, Авито Недвижимость, Авито Услуги)
    year: 2026, // Год рекапа
    locale: 'ru-RU', // Локаль – по ней на фронте ы определяем как форматировать даты – для нас всегда ru-RU
    user: { // Инфа о пользователе
      id: 'user-42',
      displayName: 'Алекс',
      avatarUrl: 'https://www.gravatar.com/avatar/?d=mp',
    },
    generatedAt: '2026-12-28T10:00:00.000Z', // Дата генерации рекапа
  },
  metrics: { // Словарь всех метрик по ключам
    listingsPublished: { type: EMetricType.Number, value: 47 },
    listingsPercentile: { type: EMetricType.Percentile, value: 88 },
    viewsTotal: { type: EMetricType.Number, value: 12840 },
    viewsPercentile: { type: EMetricType.Percentile, value: 92 },
    favoritesReceived: { type: EMetricType.Number, value: 316 },
    favoritesPercentile: { type: EMetricType.Percentile, value: 79 },
    messagesSent: { type: EMetricType.Number, value: 892 },
    messagesPercentile: { type: EMetricType.Percentile, value: 85 },
    dealsClosed: { type: EMetricType.Number, value: 23 },
    dealsPercentile: { type: EMetricType.Percentile, value: 74 },
    moneyEarned: { type: EMetricType.Money, value: 186500, currency: 'RUB' },
    moneySaved: { type: EMetricType.Money, value: 24300, currency: 'RUB' },
    deliveryOrders: { type: EMetricType.Number, value: 18 },
    searchQueries: { type: EMetricType.Number, value: 1204 },
    daysActive: { type: EMetricType.Number, value: 214 },
    topCategoryShare: { type: EMetricType.Ratio, value: 0.41 },
    peakDayViews: { type: EMetricType.Number, value: 1840 },
    categoriesTried: { type: EMetricType.Number, value: 9 },
    topCategory: { type: EMetricType.String, value: 'Электроника' },
    peakMonth: { type: EMetricType.String, value: 'сентябрь' },
    favoriteCategories: {
      type: EMetricType.List,
      value: [
        { id: 'electronics', label: 'Электроника', value: 19 },
        { id: 'home', label: 'Для дома', value: 11 },
        { id: 'hobby', label: 'Хобби', value: 8 },
      ],
    },
  },
  badges: [ // Словарь всех бейджей по ключам
    {
      id: 'early-bird',
      title: 'Ранняя пташка',
      description: 'Чаще всего публиковали объявления до 9 утра',
      iconUrl: '/icon1.svg',
    },
    {
      id: 'deal-maker',
      title: 'Мастер сделок',
      description: 'Закрыли больше 20 сделок за год',
      iconUrl: '/icon2.svg',
    },
    {
      id: 'explorer',
      title: 'Исследователь',
      description: 'Заглянули в 9 разных категорий',
      iconUrl: '/icon3.svg',
    },
  ],
  features: { // Фича флаги (пока на фронте не работают) – например возможность поделиться, 
    shareEnabled: true,
    upsellId: 'delivery-plus',
  },
  story: [ // Самое важное – массив со сторями, подробнее ниже
    {
      id: 'intro',
      type: ESceneType.Intro,
      motion: EMotionPreset.Fade,
      title: 'Алекс, ваш 2026 на Авито',
      subtitle: 'Год находок, сделок и удачных объявлений',
      actions: [{ type: ESceneActionType.Next, label: 'Смотреть итоги' }],
    },
    {
      id: 'stat-listings',
      type: ESceneType.Stat,
      value: 'listingsPublished',
      percentile: 'listingsPercentile',
      unit: { one: 'объявление', few: 'объявления', many: 'объявлений' },
      title: 'вы опубликовали',
      eyebrow: 'За 2026 год',
      blockMotion: EMotionPreset.CountUp,
      motion: EMotionPreset.SlideUp,
    },
    {
      id: 'stat-views',
      type: ESceneType.Stat,
      value: 'viewsTotal',
      percentile: 'viewsPercentile',
      unit: { one: 'просмотр', few: 'просмотра', many: 'просмотров' },
      title: 'собрали ваши объявления',
      comparisonTemplate: 'это больше, чем у {{percentile}}% пользователей',
      blockMotion: EMotionPreset.CountUp,
    },
    {
      id: 'insight-category',
      type: ESceneType.Insight,
      title: 'Ваша стихия',
      text: 'Электроника — ваша главная категория. Здесь вы были особенно заметны.',
      blockMotion: EMotionPreset.StaggerText,
      motion: EMotionPreset.SlideLeft,
      linksTo: 'stat-favorites',
    },
    {
      id: 'stat-favorites',
      type: ESceneType.Stat,
      value: 'favoritesReceived',
      percentile: 'favoritesPercentile',
      unit: { one: 'добавление', few: 'добавления', many: 'добавлений' },
      title: 'в избранное',
      eyebrow: 'Любимчики покупателей',
    },
    {
      id: 'achievement-early-bird',
      type: ESceneType.Achievement,
      badgeId: 'early-bird',
      blockMotion: EMotionPreset.BadgePop,
      motion: EMotionPreset.ScaleFade,
    },
    {
      id: 'stat-messages',
      type: ESceneType.Stat,
      value: 'messagesSent',
      percentile: 'messagesPercentile',
      unit: { one: 'сообщение', few: 'сообщения', many: 'сообщений' },
      title: 'в чатах с покупателями',
      eyebrow: 'Диалоги',
    },
    {
      id: 'insight-style',
      type: ESceneType.Insight,
      text: 'Вы отвечали быстро и по делу — покупатели это ценят.',
      blockMotion: EMotionPreset.StaggerText,
    },
    {
      id: 'stat-deals',
      type: ESceneType.Stat,
      value: 'dealsClosed',
      percentile: 'dealsPercentile',
      unit: { one: 'сделка', few: 'сделки', many: 'сделок' },
      title: 'успешно закрыто',
      eyebrow: 'Результат',
      valueFormat: { maximumFractionDigits: 0 },
    },
    {
      id: 'achievement-deal-maker',
      type: ESceneType.Achievement,
      badgeId: 'deal-maker',
      blockMotion: EMotionPreset.BadgePop,
    },
    {
      id: 'stat-earned',
      type: ESceneType.Stat,
      value: 'moneyEarned',
      unit: '₽',
      title: 'заработали на продажах',
      eyebrow: 'Доход',
      valueFormat: { maximumFractionDigits: 0 },
      blockMotion: EMotionPreset.CountUp,
    },
    {
      id: 'blocks-year-mix',
      type: ESceneType.Blocks,
      motion: EMotionPreset.SlideUp,
      blocks: [
        {
          type: ESceneBlockType.Stat,
          value: 'daysActive',
          unit: { one: 'день', few: 'дня', many: 'дней' },
          title: 'были активны',
          eyebrow: 'Присутствие',
          blockMotion: EMotionPreset.CountUp,
        },
        {
          type: ESceneBlockType.Text,
          text: 'Вы заглядывали на Авито чаще, чем в выходные к друзьям.',
          blockMotion: EMotionPreset.StaggerText,
        },
        {
          type: ESceneBlockType.Callout,
          text: '214 активных дней — почти каждый день что-то происходило',
          blockMotion: EMotionPreset.CalloutIn,
        },
      ],
    },
    {
      id: 'insight-peak',
      type: ESceneType.Insight,
      title: 'Пиковый день',
      text: 'В сентябре одно объявление набрало рекорд просмотров — настоящий хит сезона.',
      blockMotion: EMotionPreset.StaggerText,
      actions: [
        { type: ESceneActionType.Next, label: 'Дальше' },
        {
          type: ESceneActionType.GoTo,
          label: 'К началу',
          sceneId: 'intro',
          variant: EButtonVariant.Ghost,
        },
      ],
    },
    {
      id: 'stat-peak-views',
      type: ESceneType.Stat,
      value: 'peakDayViews',
      unit: { one: 'просмотр', few: 'просмотра', many: 'просмотров' },
      title: 'за один день',
      eyebrow: 'Рекорд',
      blockMotion: EMotionPreset.CountUp,
      motion: { enter: EMotionPreset.ScaleFade, durationMs: 520 },
    },
    {
      id: 'achievement-explorer',
      type: ESceneType.Achievement,
      badgeId: 'explorer',
      title: 'Исследователь категорий',
      description: '9 категорий — широкий кругозор продавца',
      blockMotion: EMotionPreset.BadgePop,
    },
    {
      id: 'upsell-delivery',
      type: ESceneType.Upsell,
      title: 'Авито Доставка',
      text: 'С доставкой вы сэкономили около {{value}} на логистике и встречи.',
      callout: 'В следующем году можно сэкономить ещё больше',
      value: 'moneySaved',
      blockMotion: EMotionPreset.CalloutIn,
      actions: [
        {
          type: ESceneActionType.Link,
          label: 'Подключить доставку',
          href: 'https://www.avito.ru/',
          variant: EButtonVariant.Primary,
        },
        { type: ESceneActionType.Next, label: 'Пропустить', variant: EButtonVariant.Ghost },
      ],
    },
    {
      id: 'blocks-search',
      type: ESceneType.Blocks,
      blocks: [
        {
          type: ESceneBlockType.Stat,
          value: 'searchQueries',
          unit: { one: 'поиск', few: 'поиска', many: 'поисков' },
          title: 'запросов за год',
          eyebrow: 'Любопытство',
          blockMotion: EMotionPreset.CountUp,
        },
        {
          type: ESceneBlockType.Callout,
          text: 'Вы не только продавали — вы и сами охотились за находками',
          blockMotion: EMotionPreset.CalloutIn,
        },
      ],
    },
    {
      id: 'custom-top-categories',
      type: ESceneType.Custom,
      sceneType: 'top-categories',
      props: {
        metricKey: 'favoriteCategories',
        title: 'Топ категорий',
      },
      motion: EMotionPreset.SlideUp,
      actions: [{ type: ESceneActionType.Next, label: 'Дальше' }],
    },
    {
      id: 'upsell-promo',
      type: ESceneType.Upsell,
      title: 'Продвижение',
      text: 'В 2027 ваши объявления могут найти ещё больше покупателей.',
      callout: 'Попробуйте пакет продвижения на старте сезона',
      blockMotion: EMotionPreset.CalloutIn,
      actions: [
        {
          type: ESceneActionType.Link,
          label: 'Узнать о продвижении',
          href: 'https://www.avito.ru/',
          variant: EButtonVariant.Secondary,
        },
        { type: ESceneActionType.Next, label: 'К финалу' },
      ],
    },
    {
      id: 'outro',
      type: ESceneType.Outro,
      title: 'Это был ваш год на Авито',
      subtitle: 'Сохраните итоги или вернитесь к объявлениям',
      motion: EMotionPreset.Fade,
      actions: [
        {
          type: ESceneActionType.Share,
          label: 'Поделиться',
          share: {
            kind: EShareKind.Link,
            title: 'Мои итоги 2026 на Авито',
            text: 'Посмотрите, каким был мой год на Авито!',
            url: 'https://www.avito.ru/',
          },
        },
        {
          type: ESceneActionType.Custom,
          id: 'close-recap',
          label: 'На главную',
          variant: EButtonVariant.Primary,
        },
      ],
    },
  ],
};
```

## Метрики

Метрики бывают 6-ти типов:

```ts
export enum EMetricType {
  Number = 'number', // Обычное число
  Money = 'money', // Денежная сумма
  Percentile = 'percentile', // Процент
  Ratio = 'ratio', // Доля
  String = 'string', // Строка
  List = 'list', // Список
}

export type MetricListItem = {
  id: string;
  label: string;
  value?: number;
  imageUrl?: string;
};

type NumberMetricValue = { type: EMetricType.Number; value: number; unit?: string };

type MoneyMetricValue = { type: EMetricType.Money; value: number; currency: string };

type PercentileMetricValue = { type: EMetricType.Percentile; value: number };

type RatioMetricValue = { type: EMetricType.Ratio; value: number };

type StringMetricValue = { type: EMetricType.String; value: string };

type ListMetricValue = { type: EMetricType.List; value: MetricListItem[] };
```

## Бейджи

Бейдж представляет из себя следующий объект:

```ts
export type Badge = {
  id: string;
  title: string;
  description: string;
  icon?: string; // URL иконки
};
```

Экраны с бейджем выглядит следующим образом:

![Экраны с бейджем](/docs/static/badge.png)

## Истории

Это самое иннтересное. Сюда надо прокидывать массив с самими сценами в порядке их отображения на экране. Библиотекой из под капота предоставляется 7 типов слайдов и один тип Custom. Кастомные сцены могут создать разработчики самостоятельно. Для нас достаточно будет этих 8-ми:

```ts
export enum ESceneType {
  Intro = 'intro', // Вступление
  Stat = 'stat', // Большая цифра
  Insight = 'insight', // Текстовый инсайт
  Achievement = 'achievement', // Ачивка
  Upsell = 'upsell', // Промо
  Blocks = 'blocks', // Несколько блоков на одной сцене
  Outro = 'outro', // Финал
  Custom = 'custom', // Кастомная сцена
}
```

Каждая сцена имеет следующие базовые параметры:

```ts
export type StoryItemBase = {
  id: string; // айди сцены – должен быть уникальным для каждой сцены
  type: ESceneType; // Тип сцены один из ESceneType
  motion?: EMotionPreset | MotionConfig; // Анимация переходов сцены – подробнее ниже
  durationMs?: number; // Длительность в милисекундах при автоплее
  actions?: SceneAction[]; // Кнопки внизу каждой страницы
};
```


По каждой сцене по-подробнее:

### Интро (intro)

![Интро](/docs/static/intro.png)

```ts
export type StoryIntroItem = StoryItemBase & {
  type: ESceneType.Intro;
  title?: string; // Большая надпись "Алекс, ваш 2026 год на Авито"
  subtitle?: string; // Надпись по-меньше "Год находок, сделок и удачных объявлений"
};
```

### Большая цифра (stat)

![Интро](/docs/static/stat.png)

```ts
export type StoryStatItem = StoryItemBase & {
  type: ESceneType.Stat;
  value: string; // Ключ числовой метрики – в примере "listingsPublished" – 47
  percentile?: string; // (Опционально) Ключ процентной метрики – в примере "listingsPercentile" – 88%
  unit?: string | PluralForms; // (Опционально) unit – строка (например "шт.") или объект по падежам (смотри ниже)
  title?: string; // (Опционально) Большой текст – в примере "вы опубликовали"
  eyebrow?: string; // (Опционально) маленький текст над title – в примере "вы опубликовали"
  comparisonTemplate?: string; // (Опционально) Строковый шаблон для сравнения, в примере это "это больше, чем у {{percentile}}% пользователей", percentile берется из метрики, которую вы передали в "listingsPercentile"
  valueFormat?: Intl.NumberFormatOptions; // (опционально) как форматировать большое число – см. ниже
  blockMotion?: EMotionPreset.CountUp | EMotionPreset.None; // Анимация самой цифры (count-up / без анимации)
};

export type PluralForms = { // Можете загуглить че такое плюральные формы
  one: string; // единственное число (например, "1 объявление")
  few: string; // от 2 до 4 штук или другие числа, оканчивающиеся на 2, 3, 4 (например, "24 объявления")
  many: string; // 0, 5–20 и т.д. (например, "15 объявлений", "47 объявлений")
};
```

`valueFormat` — это опции стандартного `Intl.NumberFormat` (локаль берётся из `meta.locale`). Форматирует **только большую цифру**, не `unit`.

В примере для денег:

```ts
{
  id: 'stat-earned',
  type: ESceneType.Stat,
  value: 'moneyEarned', // 186500
  unit: '₽',            // подпись рядом с числом
  valueFormat: { maximumFractionDigits: 0 }, // без копеек → "186 500"
}
```

Частые опции: `maximumFractionDigits`, `minimumFractionDigits`, `useGrouping`, `notation: "compact"`.  
Если поставить `style: "currency"` + `currency: "RUB"`, знак ₽ появится **внутри** числа — тогда `unit: '₽'` лучше не дублировать.

---

### Текстовый инсайт (insight)

![Инсайт](/docs/static/insight.png)

```ts
export type StoryInsightItem = StoryItemBase & {
  type: ESceneType.Insight;
  title?: string; // (опционально) заголовок – в примере "Ваша стихия"
  text?: string; // основной текст – в примере "Электроника — ваша главная категория..."
  linksTo?: string; // (опционально) id другой сцены – смысловая связь, в примере "stat-favorites"
  blockMotion?: EMotionPreset.StaggerText | EMotionPreset.None; // анимация появления текста
};
```

Пример из mock:

```ts
{
  id: 'insight-category',
  type: ESceneType.Insight,
  title: 'Ваша стихия',
  text: 'Электроника — ваша главная категория. Здесь вы были особенно заметны.',
  blockMotion: EMotionPreset.StaggerText,
  motion: EMotionPreset.SlideLeft,
  linksTo: 'stat-favorites',
}
```

Можно и без `title` — только текст, как в `insight-style`.

---

### Ачивка (achievement)

![Ачивка](/docs/static/achievement.png)

```ts
export type StoryAchievementItem = StoryItemBase & {
  type: ESceneType.Achievement;
  badgeId?: string; // (опционально) id из массива badges – подтянет title / description / icon
  title?: string; // (опционально) свой заголовок (перебивает badge.title)
  description?: string; // (опционально) своё описание (перебивает badge.description)
  icon?: string; // (опционально) свой URL иконки (перебивает badge.icon)
  blockMotion?: EMotionPreset.BadgePop | EMotionPreset.None;
};
```

Два рабочих варианта:

1. Только ссылка на бейдж (как `achievement-early-bird`):

```ts
{
  id: 'achievement-early-bird',
  type: ESceneType.Achievement,
  badgeId: 'early-bird', // → title/description/iconUrl из badges
  blockMotion: EMotionPreset.BadgePop,
}
```

2. Бейдж + свои тексты поверх (как `achievement-explorer`):

```ts
{
  id: 'achievement-explorer',
  type: ESceneType.Achievement,
  badgeId: 'explorer',
  title: 'Исследователь категорий', // перебивает title из badges
  description: '9 категорий — широкий кругозор продавца',
  blockMotion: EMotionPreset.BadgePop,
}
```

Приоритет: поля сцены → данные из `badges` по `badgeId`.

---

### Промо (upsell)

![Upsell](/docs/static/upsell.png)

```ts
export type StoryUpsellItem = StoryItemBase & {
  type: ESceneType.Upsell;
  title?: string; // заголовок – в примере "Авито Доставка"
  text?: string; // основной текст; можно вставить {{value}}
  callout?: string; // выделенная плашка
  value?: string; // (опционально) ключ money-метрики для подстановки {{value}}
  blockMotion?: EMotionPreset.CalloutIn | EMotionPreset.None;
};
```

Если в `text` / `callout` есть `{{value}}`, а в `value` передан ключ метрики типа `money`, фронт подставит отформатированную сумму (с валютой).

Пример из mock:

```ts
{
  id: 'upsell-delivery',
  type: ESceneType.Upsell,
  title: 'Авито Доставка',
  text: 'С доставкой вы сэкономили около {{value}} на логистике и встречи.',
  // {{value}} → "24 300 ₽" из metrics.moneySaved
  callout: 'В следующем году можно сэкономить ещё больше',
  value: 'moneySaved',
  blockMotion: EMotionPreset.CalloutIn,
  actions: [
    {
      type: ESceneActionType.Link,
      label: 'Подключить доставку',
      href: 'https://www.avito.ru/',
      variant: EButtonVariant.Primary,
    },
    { type: ESceneActionType.Next, label: 'Пропустить', variant: EButtonVariant.Ghost },
  ],
}
```

Upsell без подстановки суммы — обычные строки, как в `upsell-promo`.

---

### Несколько блоков (blocks)

![Blocks](/docs/static/blocks.png)

Одна сцена = несколько кусков контента друг под другом.

```ts
export enum ESceneBlockType {
  Stat = 'stat',       // маленькая цифра (как Stat, но внутри блока)
  Text = 'text',       // текст
  Callout = 'callout', // выделенный callout
}

export type StoryStatBlock = {
  type: ESceneBlockType.Stat;
  value: string;
  percentile?: string;
  unit?: string | PluralForms;
  title?: string;
  eyebrow?: string;
  comparisonTemplate?: string;
  valueFormat?: Intl.NumberFormatOptions;
  blockMotion?: EMotionPreset.CountUp | EMotionPreset.None;
};

export type StoryTextBlock = {
  type: ESceneBlockType.Text;
  text: string;
  blockMotion?: EMotionPreset.StaggerText | EMotionPreset.None;
};

export type StoryCalloutBlock = {
  type: ESceneBlockType.Callout;
  text: string;
  blockMotion?: EMotionPreset.CalloutIn | EMotionPreset.None;
};

export type StoryBlock = StoryStatBlock | StoryTextBlock | StoryCalloutBlock;

export type StoryBlocksItem = StoryItemBase & {
  type: ESceneType.Blocks;
  blocks: StoryBlock[]; // массив блоков в порядке отрисовки
};
```

Пример из mock (`blocks-year-mix`):

```ts
{
  id: 'blocks-year-mix',
  type: ESceneType.Blocks,
  motion: EMotionPreset.SlideUp,
  blocks: [
    {
      type: ESceneBlockType.Stat,
      value: 'daysActive', // 214
      unit: { one: 'день', few: 'дня', many: 'дней' },
      title: 'были активны',
      eyebrow: 'Присутствие',
      blockMotion: EMotionPreset.CountUp,
    },
    {
      type: ESceneBlockType.Text,
      text: 'Вы заглядывали на Авито чаще, чем в выходные к друзьям.',
      blockMotion: EMotionPreset.StaggerText,
    },
    {
      type: ESceneBlockType.Callout,
      text: '214 активных дней — почти каждый день что-то происходило',
      blockMotion: EMotionPreset.CalloutIn,
    },
  ],
}
```

---

### Финал (outro)

![Outro](/docs/static/outro.png)

```ts
export type StoryOutroItem = StoryItemBase & {
  type: ESceneType.Outro;
  title?: string; // большая надпись – в примере "Это был ваш год на Авито"
  subtitle?: string; // текст поменьше
};
```

Обычно на outro вешают `actions`: поделиться / закрыть / на главную.

Пример из mock:

```ts
{
  id: 'outro',
  type: ESceneType.Outro,
  title: 'Это был ваш год на Авито',
  subtitle: 'Сохраните итоги или вернитесь к объявлениям',
  motion: EMotionPreset.Fade,
  actions: [
    {
      type: ESceneActionType.Share,
      label: 'Поделиться',
      share: {
        kind: EShareKind.Link,
        title: 'Мои итоги 2026 на Авито',
        text: 'Посмотрите, каким был мой год на Авито!',
        url: 'https://www.avito.ru/',
      },
    },
    {
      type: ESceneActionType.Custom,
      id: 'close-recap', // фронт сам решит, что делать по этому id
      label: 'На главную',
      variant: EButtonVariant.Primary,
    },
  ],
}
```

---

### Кастомная сцена (custom)

![Custom](/docs/static/custom.png)


Если не хватает существующих пресетов, то можно создать кастомную сцену. В таком случае фронту с бэком нужно будет согласовать 

---

## Кнопки (actions)

Массив `actions` на любой сцене — кнопки внизу. Если не передать — фронт рисует дефолтные Prev/Next (кроме intro, там дефолт «Начать»).

```ts
export enum ESceneActionType {
  Next = 'next',     // следующая сцена
  Prev = 'prev',     // предыдущая сцена
  Link = 'link',     // переход по URL
  Share = 'share',   // шаринг
  GoTo = 'goto',     // прыжок на сцену по id
  Custom = 'custom', // кастомное событие на фронт
}

export enum EButtonVariant {
  Primary = 'primary',
  Secondary = 'secondary',
  Ghost = 'ghost',
}

export enum EShareKind {
  Story = 'story',
  Link = 'link',
}

export enum ELinkTarget {
  Blank = '_blank',
  Self = '_self',
}

export type SceneAction =
  | { type: ESceneActionType.Next; label?: string; variant?: EButtonVariant }
  | { type: ESceneActionType.Prev; label?: string; variant?: EButtonVariant }
  | {
      type: ESceneActionType.Link;
      label: string;
      href: string;
      target?: ELinkTarget;
      variant?: EButtonVariant;
    }
  | {
      type: ESceneActionType.Share;
      label: string;
      share: { kind: EShareKind; title?: string; text?: string; url?: string };
      variant?: EButtonVariant;
    }
  | {
      type: ESceneActionType.GoTo;
      label: string;
      sceneId: string; // id сцены из story
      variant?: EButtonVariant;
    }
  | {
      type: ESceneActionType.Custom;
      label: string;
      id: string; // произвольный id события для фронта
      variant?: EButtonVariant;
    };
```

Примеры из mock:

- intro → `{ type: 'next', label: 'Смотреть итоги' }`
- insight-peak → `next` + `goto` на `'intro'`
- upsell-delivery → `link` + `next` («Пропустить»)
- outro → `share` + `custom` с `id: 'close-recap'`

---

## Анимации (motion / blockMotion)

Два уровня:

1. **`motion`** на сцене — анимация **перехода** всей сцены (вход/выход).
2. **`blockMotion`** — анимация **внутри** контента (цифра count-up, текст stagger, бейдж pop, callout).

```ts
export enum EMotionPreset {
  Fade = 'fade',
  SlideUp = 'slide-up',
  SlideLeft = 'slide-left',
  ScaleFade = 'scale-fade',
  CountUp = 'count-up',         // обычно для blockMotion у Stat
  BadgePop = 'badge-pop',       // обычно для blockMotion у Achievement
  StaggerText = 'stagger-text', // обычно для Insight / Text
  CalloutIn = 'callout-in',     // обычно для Upsell / Callout
  None = 'none',
}

export type MotionConfig = {
  enter?: EMotionPreset;
  exit?: EMotionPreset;
  durationMs?: number;
  staggerMs?: number;
  ease?: EMotionEase;
};
```

`motion` можно передать строкой-пресетом или объектом:

```ts
motion: EMotionPreset.SlideUp

// или подробнее:
motion: { enter: EMotionPreset.ScaleFade, durationMs: 520 }
```

Для бэка достаточно пресетов из примера. Если не указать — фронт подставит дефолты.

---

## Narrative (опционально)

В корне payload можно добавить `narrative` — слой **текстов** поверх `story` (без изменения структуры сцен).

```ts
narrative?: {
  scenes?: Record<string, {
    title?: string;
    subtitle?: string;
    body?: string;       // для insight / upsell это основной текст
    comparison?: string; // для stat — comparisonTemplate; для upsell — callout
  }>;
  highlights?: string[]; // короткие тезисы (шаринг/превью; в плеере сейчас не рисуются)
};
```

Приоритет текста: **narrative → поля story → дефолт движка**.

Нужен, когда `story` шаблонный (одни и те же сцены/метрики), а формулировки персональные и собираются отдельно.  
Если тексты уже лежат в каждом `StoryItem` — `narrative` можно не слать (в текущем mock его нет).

---

## Media (опционально)

```ts
media?: Record<string, { url: string; alt?: string }>;
```

Справочник картинок по ключу. Задел под фоны/иллюстрации. В текущем mock не используется — можно пока игнорировать.

---

## Чеклист для бэка

Минимум, без которого виджет не собрать:

1. `schemaVersion: 1`
2. `meta` (хотя бы `year`, `locale`, `user.displayName`)
3. `metrics` — все ключи, на которые ссылается `story`
4. `story` — упорядоченный плейлист сцен
5. `badges` — если есть сцены `achievement` с `badgeId`

Опционально: `narrative`, `features`, `media`, `actions`/`motion` на сценах.
