package {{.Package}};

import lark.task.TaskApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Bootstrap {
    public static void main(String[] args) {
        TaskApplication app = new TaskApplication(Bootstrap.class);
        app.run(args);
    }
}