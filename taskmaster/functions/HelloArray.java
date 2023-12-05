import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
public class HelloArray {
    public static JsonArray main(JsonObject args) {
        JsonArray jsonArray = new JsonArray();
        jsonArray.add("a");
        jsonArray.add("b");
        return jsonArray;
    }
}