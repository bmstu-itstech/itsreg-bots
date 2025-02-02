// Package bots provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.2.0 DO NOT EDIT.
package bots

import (
	"time"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for BlockType.
const (
	Message   BlockType = "message"
	Question  BlockType = "question"
	Selection BlockType = "selection"
)

// Defines values for BotStatus.
const (
	Failed  BotStatus = "failed"
	Started BotStatus = "started"
	Stopped BotStatus = "stopped"
)

// Block Минимальная структурная единица сценария бота. Представляет из себя сообщение, которое отправляет бот пользователю, и в зависимости от типа блока обрабатывается по-разному:
//   - Сообщение (message) - просто сообщение от бота. Не ждет ответа пользователя и сразу переключает пользователя на следующий блок с состоянием next.
//   - Вопрос (question) - сообщение от бот, ожидается ответ пользователя. После ответа пользователя переключает пользователя на следующий блок с состоянием next
//   - Выбор (selection) - сообщение от бота, после которого ожидается ответ пользователя кнопкой или произвольным текстом. Если пользователь отвечает кнопкой, бот переключает его на следующий блок с состоянием next у выбранной опции (Option). Если пользователь отвечает произвольным текстом, переключает пользователя на следующий блок с состоянием next.
type Block struct {
	// NextState Состояние (state) другого блока. Конкретное значение определяется типом (type) блока.
	NextState int `json:"nextState"`

	// Options Опции для блока с выбором ответа. Не допускается использование опций для других типов блока.
	Options *[]Option `json:"options,omitempty"`

	// State Уникальный идентификатор блока в рамках бота. Не может равняться нулю.
	State int `json:"state"`

	// Text Текст сообщения бота. Не допускается использование вёрстки.
	Text string `json:"text"`

	// Title Название блока. Используется в заголовке таблицы с ответами участников.
	Title string `json:"title"`

	// Type Тип кнопки:
	//  - Сообщение (message) - просто сообщение от бота. Не ждет ответа пользователя и сразу переключает пользователя на следующий блок с состоянием next.
	//  - Вопрос (question) - сообщение от бот, ожидается ответ пользователя. После ответа пользователя переключает пользователя на следующий блок с состоянием next
	//  - Выбор (selection) - сообщение от бота, после которого ожидается ответ пользователя кнопкой или произвольным текстом. Если пользователь отвечает кнопкой, бот переключает его на следующий блок с состоянием next у выбранной опции (Option). Если пользователь отвечает произвольным текстом, переключает пользователя на следующий блок с состоянием next.
	Type BlockType `json:"type"`
}

// BlockType Тип кнопки:
//   - Сообщение (message) - просто сообщение от бота. Не ждет ответа пользователя и сразу переключает пользователя на следующий блок с состоянием next.
//   - Вопрос (question) - сообщение от бот, ожидается ответ пользователя. После ответа пользователя переключает пользователя на следующий блок с состоянием next
//   - Выбор (selection) - сообщение от бота, после которого ожидается ответ пользователя кнопкой или произвольным текстом. Если пользователь отвечает кнопкой, бот переключает его на следующий блок с состоянием next у выбранной опции (Option). Если пользователь отвечает произвольным текстом, переключает пользователя на следующий блок с состоянием next.
type BlockType string

// Bot Информация о боте.
type Bot struct {
	// Blocks Все блоки бота, см. Block.
	Blocks []Block `json:"blocks"`

	// BotUUID Уникальный идентификатор бота. Не должен превышать длину в 36 символов.
	BotUUID string `json:"botUUID"`

	// CreatedAt Время создания бота.
	CreatedAt time.Time `json:"createdAt"`

	// Entries Все точки входа бота, см. EntryPoint. Гарантировано существует точка входа start
	Entries []EntryPoint `json:"entries"`

	// Mailings Все рассылки бота, см. Mailings.
	Mailings *[]Mailing `json:"mailings,omitempty"`

	// Name Имя бота.
	Name string `json:"name"`

	// Status Статус бота: started (запущен), stopped (не запущен), failed (ошибка запуска).
	Status BotStatus `json:"status"`

	// Token Телеграм токен бота. Получить токен можно в телеграм-боте @BotFather.
	Token string `json:"token"`

	// UpdatedAt Время последнего обновления бота.
	UpdatedAt time.Time `json:"updatedAt"`
}

// BotStatus Статус бота: started (запущен), stopped (не запущен), failed (ошибка запуска).
type BotStatus string

// CreateMailing Рассылка и связанные с ней точка входа и блоки.
type CreateMailing struct {
	// Blocks Список блоков для рассылки. Обычно содержит единственный блок типа message.
	Blocks []Block `json:"blocks"`

	// EntryPoint Точка входа для бота. Бот должен иметь как минимум точку входа "start". Иные точки входа используются для создания рассылок. Точка входа начинает скрипт бота с отправки блока с состоянием state пользователю.
	EntryPoint EntryPoint `json:"entryPoint"`

	// Name Название рассылки.
	Name string `json:"name"`

	// RequiredState Состояние участника, необходимое для получения рассылки.
	RequiredState int `json:"requiredState"`
}

// EntryPoint Точка входа для бота. Бот должен иметь как минимум точку входа "start". Иные точки входа используются для создания рассылок. Точка входа начинает скрипт бота с отправки блока с состоянием state пользователю.
type EntryPoint struct {
	// Key Уникальный ключ точки входа бота.
	Key string `json:"key"`

	// State Состояние (state) первого блока в скрипте.
	State int `json:"state"`
}

// Error Описание ошибки.
type Error struct {
	Message string `json:"message"`
}

// GetBots Список ботов.
type GetBots = []Bot

// Mailing Рассылка от бота. При старте рассылки активирует точку входа с ключом entryKey всем пользователям, прошедшим блок с состоянием requiredState.
type Mailing struct {
	// EntryKey Ключ точки входа (EntryPoint), которая активируется при старте рассылки.
	EntryKey string `json:"entryKey"`

	// Name Имя рассылки.
	Name string `json:"name"`

	// RequiredState Состояние (state) блока, требуемого для отправки рассылки или 0, для всех, кто завершил скрипт. Так, requiredState = 5 обозначает, что рассылка будет отправлена всем участникам, которые ответили на блок с состоянием 5.
	RequiredState int `json:"requiredState"`
}

// Option Опция для блока с выбором ответа. Представлена в telegram как кнопка в клавиатуре (ReplyKeyboard).
type Option struct {
	// Next Состояние (state) следующего блока, если пользователь выбрал данную опцию.
	Next int `json:"next"`

	// Text Текст на кнопке.
	Text string `json:"text"`
}

// PostBots Данные, необходимые для создания бота.
type PostBots struct {
	// Blocks Все блоки бота, см. Block.
	Blocks []Block `json:"blocks"`

	// BotUUID Уникальный идентификатор бота. Не должен превышать длину в 36 символов.
	BotUUID string `json:"botUUID"`

	// Entries Все точки входа бота, см. EntryPoint. Необходимо наличие точки входа start.
	Entries []EntryPoint `json:"entries"`

	// Mailings Все рассылки бота, см. Mailings.
	Mailings *[]Mailing `json:"mailings,omitempty"`

	// Name Имя бота.
	Name string `json:"name"`

	// Token Телеграм токен бота. Получить токен можно в телеграм-боте @BotFather.
	Token string `json:"token"`
}

// CreateBotJSONRequestBody defines body for CreateBot for application/json ContentType.
type CreateBotJSONRequestBody = PostBots

// CreateMailingJSONRequestBody defines body for CreateMailing for application/json ContentType.
type CreateMailingJSONRequestBody = CreateMailing
