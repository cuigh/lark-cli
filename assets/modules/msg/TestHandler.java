package {{.Package}}.handler;

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