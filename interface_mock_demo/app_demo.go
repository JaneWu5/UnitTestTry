package interface_mock_demo

// gomock是Go官方提供的测试框架， 它在内置的testing包或其他环境中都能够很方便的使用。
// 我们使用它对代码中的那些接口类型进行mock，方便编写单元测试。
// https://github.com/golang/mock

// mockgen 命令用来为给定一个包含要mock的接口的Go源文件，生成mock类源代码
// 运行mockgen:
// mockgen 有两种操作模式：源码（source）模式和反射（reflect）模式。
// 【源码模式】根据源文件mock接口。它是通过使用 -source 标志启用。在这个模式下可能有用的其他标志是 -imports 和 -aux_files。
// eg. mockgen -source=foo.go [other options]
// 【反射模式】通过构建使用反射来理解接口的程序来mock接口。
// 它是通过传递两个非标志参数来启用的：一个导入路径和一个逗号分隔的符号列表。可以使用 ”.”引用当前路径的包。
// eg.
// mockgen database/sql/driver Conn,Driver
// # Convenient for `go:generate`.
// mockgen . Conn,Driver

// 假设有查询MySQL数据库的业务代码如下，其中DB是一个自定义的接口类型
// 我们现在要为GetFromDB函数编写单元测试代码，可是我们又不能在单元测试过程中连接真实的数据库，
// 这个时候就需要mock DB这个接口来方便进行单元测试。
// mockgen -source=db.go -destination=mocks/db_mock.go -package=mocks
// DB 数据接口
type DB interface {
	Get(key string) (int, error)
	Add(key string, value int) error
}

// GetFromDB 根据key从DB查询数据的函数
func GetFromDB(db DB, key string) int {
	if v, err := db.Get(key); err == nil {
		return v
	}
	return -1
}

// Go测试库testify目前也提供类似的mock工具—testify/mock和mockery
