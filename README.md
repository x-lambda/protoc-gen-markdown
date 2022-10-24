# protoc-gen-markdown

根据`proto`文件生成`markdown`文档
> 依赖`protoc`，请先安装`protoc3`



### 安装
```shell
go get -u github.com/x-lambda/protoc-gen-markdown
```

### 使用
```shell
protoc -I . --markdown_out=filename=foo.md:/tmp/doc_path demo.proto
```

##### `TODO`
补充`test case`，手动创建一个`CodeGeneratorRequest`对象，然后使用`proto.Marshal()`序列化成字节对象，
最后通过`buffer`包装成`io.Reader`，验证输出结果是否和预期一致。

##### 参考
`CodeGeneratorRequest`对象比较复杂，可以参考已有的文档/项目来理解
> `github.com/davyxu/pbmeta`

项目将`protoc`传递的`CodeGeneratorRequest`对象输出到文件
<br>
<br>

***

> `tips`: 重复造了一个轮子的初衷是，现有的`protoc-gen-markdown`作为一个二进制文件，不好集成到其他系统中。
> 因此，想单独提供一个`package`的方式，可以让其他项目引用。
> `generator.ReadGenRequest()`可以接受不同的源，例如可以自定义一个合法的`CodeGeneratorRequest`对象，
> 放到`buffer`中，然后再通过`ReadGenRequest()`和`Generator.Generate()`转成`CodeGeneratorResponse`对象。
> 这样如果`protoc`也提供`package`的方式，则不需要安装这些二进制，可以直接在系统中集成即可。


##### TODO
markdown锚点设置