package {{.Package}};

import lark.msg.MsgApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Bootstrap {
    public static void main(String[] args) {
        MsgApplication app = new MsgApplication(Bootstrap.class);
        app.run(args);
    }
}