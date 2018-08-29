package tpl

// RPC 服务程序启动类模板
const tplServiceBootstrap = `package {{.Package}};

import lark.net.rpc.RpcApplication;
import lark.net.rpc.server.ServerOptions;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.autoconfigure.data.mongo.MongoDataAutoConfiguration;
import org.springframework.boot.autoconfigure.mongo.MongoAutoConfiguration;

@SpringBootApplication(exclude = {MongoAutoConfiguration.class, MongoDataAutoConfiguration.class})
public class Bootstrap {
    public static void main(String[] args) {
        // default options
        ServerOptions options = new ServerOptions(":[请替换成服务端口]", "[请替换成服务描述]");
        RpcApplication app = new RpcApplication(Bootstrap.class, args, options);
        app.run();
    }
}
`

// RPC 服务程序服务实现类模板
const tplTestServiceImpl = `package {{.Package}}.impl;

import lark.net.rpc.annotation.RpcService;
import {{.Package}}.biz.TestBiz;
import {{.Package}}.dto.TestDto;
import {{.Package}}.entity.TestObject;
import {{.Package}}.iface.TestService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

@Component
@RpcService(name = "TestService", description = "测试服务")
public class TestServiceImpl implements TestService {
    @Autowired
    private TestBiz testBiz;

    @Override
    public TestDto.HelloResponse hello(TestDto.HelloRequest request) {
        TestObject object = testBiz.getObject(request.getId());
        TestDto.HelloResponse response = new TestDto.HelloResponse();
        response.setResult(object.getName());
        return response;
    }
}
`

// 测试服务单元测试
const tplTestServiceTest = `package {{.Package}};

import {{.Package}}.constant.TestType;
import {{.Package}}.dto.TestDto;
import {{.Package}}.iface.TestService;
import org.junit.Test;

import java.time.LocalDateTime;

public class TestServiceTest {
    @Autowired
    private TestService testService;

    @Test
    public void testHello() throws Exception {
        TestDto.HelloRequest request = new TestDto.HelloRequest();
        request.setId(123);
        request.setType(TestType.GOOD);
        request.setTime(LocalDateTime.now());

        TestDto.HelloResponse response = testService.hello(request);
        printJsonln(response);
    }
}
`

// RPC 程序单元测试配置文件
const tplServiceTestAppProperties = `mtime.testing.remoting=true`
