package tpl

// 消息处理程序启动类模板
const tplMsgBootstrap = `package {{.Package}};

import lark.util.msg.MsgApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.autoconfigure.data.mongo.MongoDataAutoConfiguration;
import org.springframework.boot.autoconfigure.mongo.MongoAutoConfiguration;

@SpringBootApplication
public class Bootstrap {
    public static void main(String[] args) {
        MsgApplication app = new MsgApplication(Bootstrap.class, args);
        app.run();
    }
}
`

// 消息处理程序实现类模板
const tplTestHandler = `package {{.Package}}.handler;

import lark.util.msg.AbstractHandler;
import lark.util.msg.Message;
import lark.util.msg.MsgHandler;
import {{.Package}}.biz.TestBiz;
import {{.Package}}.entity.TestObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

// todo: 更改 topic, channel, threads
@Component
@MsgHandler(topic = "test", channel = "process", threads = 4)
public class TestHandler extends AbstractHandler<Integer> {
    @Autowired
    private TestBiz testBiz;

    @Override
    protected void process(Integer msg, Message raw) {
        // todo: 添加实现
        TestObject object = testBiz.getObject(msg);
        System.out.println(object.getId());
    }
}
`

// 消息处理程序实现类模板
const tplTestHandlerTest = `package {{.Package}}.handler;

import {{.Package}}.AbstractTest;
import org.junit.Test;

public class TestHandlerTest extends AbstractTest {
    @Autowired
    private TestHandler testHandler;

    @Test
    public void testProcess() throws Exception {
        testHandler.process(123, null);
    }
}
`
