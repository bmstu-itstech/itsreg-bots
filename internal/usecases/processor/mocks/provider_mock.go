package mocks

import "github.com/zhikh23/itsreg-bots/internal/domain/module"

type ProviderMock struct {
	s map[int32]*module.Module
}

func NewProviderMock() *ProviderMock {
	return &ProviderMock{
		make(map[int32]*module.Module),
	}
}

func (p *ProviderMock) Register(m *module.Module) {
	p.s[m.Id] = m
}

func (p *ProviderMock) Module(id int32) (*module.Module, error) {
	return p.s[id], nil
}
