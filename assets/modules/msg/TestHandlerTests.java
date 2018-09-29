package {{.Package}}.handler;

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