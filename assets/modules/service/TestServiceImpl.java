package {{.Package}}.impl;

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