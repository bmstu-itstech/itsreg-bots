package bots

import (
	"fmt"
	"time"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/commonerrs"
)

type Bot struct {
	UUID string

	entryPoints map[string]EntryPoint
	blocks      map[int]Block

	Name   string
	Token  string
	Status Status

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBot(
	uuid string,
	entries []EntryPoint,
	blocks []Block,
	name string,
	token string,
) (*Bot, error) {
	if uuid == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty uuid")
	}

	if len(entries) == 0 {
		return nil, commonerrs.NewInvalidInputError("expected not empty entries")
	}

	if len(blocks) == 0 {
		return nil, commonerrs.NewInvalidInputError("expected not empty blocks")
	}

	if name == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty name")
	}

	if token == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty token")
	}

	for _, entry := range entries {
		if entry.IsZero() {
			return nil, commonerrs.NewInvalidInputError("expected not empty entry")
		}
	}

	es, err := mapEntries(entries)
	if err != nil {
		return nil, err
	}

	for _, block := range blocks {
		if block.IsZero() {
			return nil, commonerrs.NewInvalidInputError("expected not empty block")
		}
	}

	bs, err := mapBlocks(blocks)
	if err != nil {
		return nil, err
	}

	vs := vertices(bs)
	for _, entry := range entries {
		err := colorizeVertices(vs, entry.State)
		if err != nil {
			return nil, err
		}
	}

	if whiteVertexState := findWhiteVertex(vs); whiteVertexState > 0 {
		return nil, newUnusedBlockFoundError(whiteVertexState)
	}

	return &Bot{
		UUID:      uuid,
		Name:      name,
		Token:     token,
		Status:    Stopped,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		entryPoints: es,
		blocks:      bs,
	}, nil
}

func MustNewBot(
	uuid string,
	entryPoints []EntryPoint,
	blocks []Block,
	name string,
	token string,
) *Bot {
	b, err := NewBot(uuid, entryPoints, blocks, name, token)
	if err != nil {
		panic(err)
	}
	return b
}

func NewBotFromDB(
	uuid string,
	entries []EntryPoint,
	blocks []Block,
	name string,
	token string,
	status string,
	createdAt time.Time,
	updatedAt time.Time,
) (*Bot, error) {
	if uuid == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty uuid")
	}

	if blocks == nil {
		return nil, commonerrs.NewInvalidInputError("expected not empty blocks")
	}

	if name == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty name")
	}

	if token == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty token")
	}

	if status == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty status")
	}

	if createdAt.IsZero() {
		return nil, commonerrs.NewInvalidInputError("expected not empty created at timestamp")
	}

	if updatedAt.IsZero() {
		return nil, commonerrs.NewInvalidInputError("expected not empty updated at timestamp")
	}

	bs, err := mapBlocks(blocks)
	if err != nil {
		return nil, err
	}

	es, err := mapEntries(entries)
	if err != nil {
		return nil, err
	}

	st, err := NewStatusFromString(status)
	if err != nil {
		return nil, err
	}

	return &Bot{
		UUID:        uuid,
		entryPoints: es,
		blocks:      bs,
		Name:        name,
		Token:       token,
		Status:      st,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (b *Bot) Children(block Block) []Block {
	children := make([]Block, 0)

	if block.NextState != 0 {
		next := b.blocks[block.NextState]
		children = append(children, next)
		children = append(children, b.Children(next)...)
	}

	for _, opt := range block.Options {
		next := b.blocks[opt.Next]
		children = append(children, next)
		children = append(children, b.Children(next)...)
	}

	return children
}

func (b *Bot) Blocks() []Block {
	blocks := make([]Block, 0, len(b.blocks))
	for _, block := range b.blocks {
		blocks = append(blocks, block)
	}
	return blocks
}

func (b *Bot) Entries() []EntryPoint {
	entries := make([]EntryPoint, 0, len(b.entryPoints))
	for _, entry := range b.entryPoints {
		entries = append(entries, entry)
	}
	return entries
}

func (b *Bot) SetBlocks(entries []EntryPoint, blocks []Block) error {
	for _, entry := range entries {
		if entry.IsZero() {
			return commonerrs.NewInvalidInputError("expected not empty entry")
		}
	}

	es, err := mapEntries(entries)
	if err != nil {
		return err
	}

	for _, block := range blocks {
		if block.IsZero() {
			return commonerrs.NewInvalidInputError("expected not empty block")
		}
	}

	bs, err := mapBlocks(blocks)
	if err != nil {
		return err
	}

	vs := vertices(bs)
	for _, entry := range entries {
		err := colorizeVertices(vs, entry.State)
		if err != nil {
			return err
		}
	}

	if whiteVertexState := findWhiteVertex(vs); whiteVertexState > 0 {
		return newUnusedBlockFoundError(whiteVertexState)
	}

	b.blocks = bs
	b.entryPoints = es
	b.UpdatedAt = time.Now()

	return nil
}

func (b *Bot) SetStatus(status Status) {
	b.Status = status
}

type vertex struct {
	Block Block
	Color color
}

type color int

const (
	white color = iota
	grey
	black
)

func newBlockNotFoundError(state int) error {
	return commonerrs.NewInvalidInputError(
		fmt.Sprintf("block with state %d not found", state),
	)
}

func newBlockIsDuplicatedError(state int) error {
	return commonerrs.NewInvalidInputError(
		fmt.Sprintf("block with state %d is duplicated", state),
	)
}

func newUnusedBlockFoundError(state int) error {
	return commonerrs.NewInvalidInputError(
		fmt.Sprintf("block with state %d is unused", state),
	)
}

func newEntryIsDuplicatedError(key string) error {
	return commonerrs.NewInvalidInputError(
		fmt.Sprintf("entry with key '%s' is duplicated", key),
	)
}

func mapBlocks(blocks []Block) (map[int]Block, error) {
	mapped := make(map[int]Block)
	for _, block := range blocks {
		if _, ok := mapped[block.State]; ok {
			return nil, newBlockIsDuplicatedError(block.State)
		}
		mapped[block.State] = block
	}
	return mapped, nil
}

func mapEntries(entries []EntryPoint) (map[string]EntryPoint, error) {
	mapped := make(map[string]EntryPoint)
	for _, entry := range entries {
		if _, ok := mapped[entry.Key]; ok {
			return nil, newEntryIsDuplicatedError(entry.Key)
		}
		mapped[entry.Key] = entry
	}
	return mapped, nil
}

func vertices(blocks map[int]Block) map[int]*vertex {
	v := make(map[int]*vertex)
	for state, block := range blocks {
		v[state] = &vertex{
			Block: block,
			Color: white,
		}
	}
	return v
}

func colorizeVertices(vertices map[int]*vertex, currentState int) error {
	current, ok := vertices[currentState]
	if !ok {
		return newBlockNotFoundError(currentState)
	}

	current.Color = grey

	childrenStates := current.Block.ChildrenStates()

	for _, nextState := range childrenStates {
		next, ok := vertices[nextState]
		if !ok {
			return newBlockNotFoundError(nextState)
		}

		if next.Color == white {
			err := colorizeVertices(vertices, nextState)
			if err != nil {
				return err
			}
			next.Color = black
		}
	}

	return nil
}

func findWhiteVertex(vertices map[int]*vertex) int {
	for _, v := range vertices {
		if v.Color == white {
			return v.Block.State
		}
	}
	return 0
}

func (b *Bot) Traverse(startState int) []Block {
	vs := vertices(b.blocks)
	b.traverseRecursive(vs, startState)

	traveled := make([]Block, 0, len(vs))
	for _, v := range filterVertices(vs, notWhitePredicate) {
		traveled = append(traveled, v.Block)
	}
	return traveled
}

func (b *Bot) traverseRecursive(vertices map[int]*vertex, currentState int) {
	current := vertices[currentState]
	current.Color = grey

	childrenStates := current.Block.ChildrenStates()

	for _, nextState := range childrenStates {
		next := vertices[nextState]

		if next.Color == white {
			b.traverseRecursive(vertices, nextState)
			next.Color = black
		}
	}
}

func filterVertices(vertices map[int]*vertex, predicate func(v *vertex) bool) []*vertex {
	res := make([]*vertex, 0)
	for _, v := range vertices {
		if predicate(v) {
			res = append(res, v)
		}
	}
	return res
}

func notWhitePredicate(v *vertex) bool {
	return v.Color != white
}