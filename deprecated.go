package xorm

// all content in this file will be removed from xorm some times after

// @deprecation : please use NewSession instead
func (engine *Engine) MakeSession() (Session, error) {
	s := engine.NewSession()
	return *s, nil
}

// @deprecation : please use NewEngine instead
func Create(driverName string, dataSourceName string) Engine {
	engine, _ := NewEngine(driverName, dataSourceName)
	return *engine
}