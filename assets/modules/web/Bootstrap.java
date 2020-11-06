package {{.Package}};

import lark.web.WebApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Bootstrap {
    public static void main(String[] args) {
        WebApplication app = new WebApplication(Bootstrap.class);
        app.run(args);
    }
}