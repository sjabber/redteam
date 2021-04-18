package com.hanium.mer.controller;

import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;

@Slf4j
@Controller
public class warningController {

    @GetMapping("/warning")
    public String redirectWarning(int tNo, int pNo){

        BufferedReader in = null;

        try {
            URL obj = new URL("http://mert.koreacentral.cloudapp.azure.com:5000/api/CountTarget?tNo="+tNo+"&pNo="+pNo+"&email=false&link=true&download=false"); // 호출할 url
            HttpURLConnection con = (HttpURLConnection)obj.openConnection();

            con.setRequestMethod("GET");
            in = new BufferedReader(new InputStreamReader(con.getInputStream(), "UTF-8"));

            String line;
            while((line = in.readLine()) != null) { // response를 차례대로 출력
                System.out.println(line);
            }
        } catch(Exception e) {
            log.error(e.getMessage());
        } finally {
            if(in != null) try { in.close(); } catch(Exception e) { e.printStackTrace(); }
        }

        return "redirect:http://mert.koreacentral.cloudapp.azure.com:8080/warn/warning2";
    }
}
