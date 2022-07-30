- gomonkey是一款go语言的打桩框架，目标是让用户在单元测试中低成本的完成打桩，从而将精力聚焦与业务功能的开发。
- gomonkey基础特性列表如下：
  - 支持一个函数打一个桩
  - 支持为一个成员方法打一个桩
  - 支持为一个全局变量打一个桩
  - 支持为一个函数变量打一个桩
  - 支持为一个函数打一个特定的桩序列
  - 支持为一个成员方法打一个特定的桩序列
  - 支持为一个函数变量打一个特定的桩序列
###### interface惯用法刷新
- 刷新1：当为interface打一个桩时，用户直接复用组合之前的ApplyFunc和ApplyMethod接口即可
  - 
```go
func TestApplyInterfaceReused(t *testing.T) {
	e := &fake.Etcd{}
	
	Convey("TestApplyInterface", t, func(){
		patches := ApplyFunc(fake.NewDb, func(_ string) fake.Db{
			return e
        })
		defer patches.Reset()
		db := fake.NewDb("mysql")
		
		Convey("TestApplyInterface", t, func(){
			info := "hello interface"
			patches.ApplyMethod(e, "Retrieve", 
				func(_ *fake.Etcd, _ string) (string, error) {
					return info, nil
                })
			output, err := db.Retrieve("")
			So(err, ShouldEqual, nil)
			So(output, ShouldEqual, info)
        })
    })
}
```