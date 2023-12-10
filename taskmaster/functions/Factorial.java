import com.google.gson.JsonObject;
public class Factorial {
    public static int fact(int n) {
        if (n <= 1)
            return n;
        return n * fact(n-1);
    }
    
    public static JsonObject main(JsonObject args) {
        int n = 0;
        if (args.has("n"))
            n = args.getAsJsonPrimitive("n").getAsInt();
        JsonObject response = new JsonObject();
        response.addProperty("result", fact(n));
        return response;
    }
}