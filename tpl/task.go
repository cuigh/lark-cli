package tpl

// 计划任务程序启动类模板
const tplTaskBootstrap = `package {{.Package}};

import lark.net.rpc.server.ServerOptions;
import lark.task.TaskApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.autoconfigure.data.mongo.MongoDataAutoConfiguration;
import org.springframework.boot.autoconfigure.mongo.MongoAutoConfiguration;

@SpringBootApplication
public class Bootstrap {
    public static void main(String[] args) {
        // default options
        ServerOptions options = new ServerOptions(":[请替换成服务端口]", "[请替换成服务描述]");        
        TaskApplication app = new TaskApplication(Bootstrap.class, args, options);
        app.run();
    }
}
`

// 计划任务程序实现类模板
const tplTestTask = `package {{.Package}}.executor;

import lark.task.Executor;
import lark.task.Task;
import lark.task.TaskContext;
import {{.Package}}.biz.TestBiz;
import {{.Package}}.entity.TestObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

@Component
@Task
public class TestTask implements Executor {
    private static final Logger LOGGER= LoggerFactory.getLogger(TestTask.class);

    @Autowired
    private TestBiz testBiz;

    @Override
    public void execute(TaskContext ctx) {
        // todo: 添加实现
        int id = ctx.getArgs().getInt32("id", -1);
        TestObject object = testBiz.getObject(id);
        System.out.println(object.getId());
    }
}
`

// 计划任务程序实现类模板
const tplTestTaskTest = `package {{.Package}}.executor;

import lark.task.TaskContext;
import lark.task.data.ExecuteParam;
import lark.util.ioc.ServiceLocator;
import {{.Package}}.AbstractTest;
import org.junit.Test;

import java.util.ArrayList;

public class TestTaskTest {
    @Autowired
    private TestTask testTask;

    @Test
    public void testExecute() throws Exception {
        ExecuteParam param = new ExecuteParam();
        param.Args = new ArrayList<>();
        ExecuteParam.Arg arg = new ExecuteParam.Arg();
        arg.Name = "id";
        arg.Value = "123";
        param.Args.add(arg);

        TaskContext ctx = new TaskContext(param);
        testTask.execute(ctx);
    }
}
`
