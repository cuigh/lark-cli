package {{.Package}}.executor;

import lark.task.TaskContext;
import lark.task.data.ExecuteParam;
import lark.util.ioc.ServiceLocator;
import {{.Package}}.AbstractTest;
import org.junit.Test;

import java.util.ArrayList;

public class TestTaskTests {
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