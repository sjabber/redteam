package com.hanium.mer.controller;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;

@Controller
public class warningController {

    @GetMapping("/warning")
    public String redirectWarning(int tNo, int pNo){

        BufferedReader in = null;

        try {
            URL obj = new URL("http://localhost:5000/api/CountTarget?tNo="+tNo+"&pNo="+pNo+"&email=false&link=true&download=false"); // 호출할 url
            HttpURLConnection con = (HttpURLConnection)obj.openConnection();

            con.setRequestMethod("GET");
            //con.setRequestProperty("Cookie","access-token="+"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MTk5OTk5OTksInVzZXJfZW1haWwiOiJnYzUyMTVAbmF0ZS5jb20iLCJ1c2VyX25hbWUiOiLquYDqtJHssYQiLCJ1c2VyX25vIjoyfQ.Dyh2Fgn_uh92s3Bv3HMV0b5fMSKBr0xfWE_h6ZBI8NA");
            in = new BufferedReader(new InputStreamReader(con.getInputStream(), "UTF-8"));

            String line;
            while((line = in.readLine()) != null) { // response를 차례대로 출력
                System.out.println(line);
            }
        } catch(Exception e) {
            e.printStackTrace();
        } finally {
            if(in != null) try { in.close(); } catch(Exception e) { e.printStackTrace(); }
        }

        return "redirect:http://localhost:8888/warn/warning2";
    }
}
