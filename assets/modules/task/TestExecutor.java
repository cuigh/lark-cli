package {{.Package}}.executor;

import lark.task.Executor;
import lark.task.Task;
import lark.task.TaskContext;
import {{.Package}}.biz.TestBiz;
import {{.Package}}.entity.TestObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

@Component
@Task
public class TestExecutor implements Executor {
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