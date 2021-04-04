package com.hanium.mer.service;

import com.hanium.mer.repogitory.SMTPRepository;
import com.hanium.mer.vo.SmtpVo;
import com.hanium.mer.vo.TargetVo;
import com.hanium.mer.vo.TemplateVO;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.compress.archivers.zip.ZipArchiveEntry;
import org.apache.commons.compress.archivers.zip.ZipArchiveOutputStream;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.FileSystemResource;
import org.springframework.mail.javamail.JavaMailSenderImpl;
import org.springframework.mail.javamail.MimeMessageHelper;
import org.springframework.stereotype.Service;

import javax.mail.internet.InternetAddress;
import javax.mail.internet.MimeMessage;
import java.io.*;
import java.util.Optional;
import java.util.Properties;
import java.util.zip.ZipEntry;
import java.util.zip.ZipOutputStream;

@Slf4j
@Service
public class SMTPService {

    @Autowired
    SMTPRepository smtpRepository;

    @Autowired
    AESService aesService;

    public static String filePath = "C:\\mailTemp\\"; // "C:\\Users\\JAVIS\\Desktop\\mer_consumer\\";

    public Optional<SmtpVo> getSMTP(Long user_no){
        Optional<SmtpVo> smtp =  smtpRepository.findByUserNo(user_no);
        try {
            //log.info(aesService.decAES(smtp.get().getSmtpPw()));//log
        }catch(Exception e){
            e.printStackTrace();
        }
        return smtp;
    }

    public String sendMail(SmtpVo smtp_info, TargetVo target, TemplateVO template, Long pNo) throws Exception {
        //TODO checkConect와 중복제거
        JavaMailSenderImpl mailSender = new JavaMailSenderImpl();
        mailSender.setHost(smtp_info.getSmtpHost());
        mailSender.setPort(Integer.parseInt(smtp_info.getSmtpPort()));

        mailSender.setUsername(smtp_info.getSmtpId());
        //비번 초기 아이디 생성시, smtp 가 자동 생성되며 비밀번호에 BCrypt 문자가 생성
        // => Illegal base64 character 24 예외
        //TODO 예외 잡을지 말지
        mailSender.setPassword(aesService.decAES(smtp_info.getSmtpPw()));

        Properties props = mailSender.getJavaMailProperties();
        props.put("mail.smtp.startttls.enable", smtp_info.getSmtpTls().equals("1")?"true":"false");//smtp_info.getSmtpTls().equals("1")?"true":"false"
        props.put("mail.smtp.auth", "true"); // 인증이 필요없다면 불필요
        props.put("mail.transport.protocol", "smtp");
        props.put("mail.smtp.connectiontimeout", smtp_info.getSmtpTimeOut());
        // 1:"smtp" 2: "smtps"
        if( smtp_info.getSmtpProtocol().equals("2") || smtp_info.getSmtpPort().equals("465")){
            //SSL을 사용하지 않으면 필요없는 것 같다. SSL 사용포트 465, TLS는 587
            props.put("mail.smtp.socketFactory.port", smtp_info.getSmtpPort());
            props.put("mail.smtp.socketFactory.class", "javax.net.ssl.SSLSocketFactory");
            props.put("mail.smtp.socketFactory.fallback", "false");
        }

        log.info("smtp info {}", smtp_info);
        //메세지 보내기: html 형태(검색시 금방 나옴), 파일도 보낼 수 있다.
        //여러 사람에게 보내기 List InternetAddress 형태? 검색필요
        //InternetAddress를 배열로 받아서 후에 프로젝트에서 target의 메일들을 받아보내면 됨
        MimeMessage message = mailSender.createMimeMessage();
        MimeMessageHelper messageHelper = new MimeMessageHelper(message, true, "UTF-8");

        //발신자
        messageHelper.setFrom(new InternetAddress(smtp_info.getSmtpId()));// + "@" + smtp_info.getSmtpHost().substring(7)));

        //수신자
        log.info("target Email {}", target.getTargetEmail());
        //message.addRecipient(Message.RecipientType.TO, new InternetAddress(target.getTargetEmail()));
        messageHelper.setTo(new InternetAddress(target.getTargetEmail()));

        // 메일 제목
        messageHelper.setSubject(template.getMailTitle());

        // 메일 내용
        messageHelper.setText(template.getMailContent(), true);
        //message.setContent(template.getMailContent(), "text/html; charset=utf-8");
        //message.setText(template.getMailContent());


        StringBuffer sb = new StringBuffer();
        sb.append(pNo).append("_").append(target.getTargetNo()).append("_bat.zip");

        //파일 첨부
        if(template.getDownloadType() == 3){
            makeBatFile(pNo + "_" + target.getTargetNo() + ".bat", target.getTargetNo(), pNo);
            makeZip(pNo + "_" + target.getTargetNo() + ".bat", pNo + "_" + target.getTargetNo() + "_temp.zip");
            makeZip(pNo + "_" + target.getTargetNo() + "_temp.zip", sb.toString());

            deleteFile(pNo + "_" + target.getTargetNo() + ".bat");
            deleteFile(pNo + "_" + target.getTargetNo() + "_temp.zip");


            FileSystemResource resource = new FileSystemResource(sb.toString());
            messageHelper.addAttachment(sb.toString(), resource);
            //messageHelper.addAttachment(sb.toString(), new ByteArrayResource(makeAttachment(target.getTargetNo(), pNo)));

        }

        /* 메모리에서 bat 파일만들고, 압축하여 보내기, 압축파일이 제대로 생성되지 않는 것 같음
        else if(template.getDownloadType() == 3) {
            byte[] zipFile = makeAttachment(target.getTargetNo(), pNo);
            if (zipFile == null){
            }else {
                ByteArrayResource resource = new ByteArrayResource(zipFile);
                messageHelper.addAttachment(sb.toString(), resource);
            }
        }
        */

        try {
            mailSender.send(message);
            if(template.getDownloadType() == 3){
                deleteFile(sb.toString());
            }
        }catch (Exception e){
            e.printStackTrace();//log.info("mail fail {}", e.printStackTrace());
            return "fail";
        }

        return "success";
    }


    //todo 각 파일라이터, 버퍼라이터 등 사용하는 이유, 장단점, 최신 트렌드, 어떤 예외들이 있는지
    public void makeBatFile(String fileName, int targetNo, Long pNo){

        try( //요기서 객체를 생성하면 try종료 후 자동으로 close처리됨
             //true : 기존 파일에 이어서 작성 (default는 false임)
             FileWriter fw = new FileWriter(fileName);
             BufferedWriter bw = new BufferedWriter(fw);
        )
        {
            bw.write("Msg * \"started ransomware contact to admin\"");
            bw.newLine();
            String countApiStr = "curl \"http://localhost:5000/api/CountTarget?tNo=" + targetNo + "&pNo=" + pNo + "&email=false&link=false&download=true\"";
            bw.write(countApiStr);
            bw.flush(); //버퍼의 내용을 파일에 쓰기
        }catch ( IOException e ) {
            log.info(e.getMessage());
        }
    }

    public void makeZip(String file4Zip, String outZipName) throws Exception{
        //아파치
        int size = 1024;
        byte[] buf = new byte[size];
        FileInputStream fis = null;
        ZipArchiveOutputStream zos = null;
        BufferedInputStream bis = null;
        try { // Zip 파일생성
            zos = new ZipArchiveOutputStream(new BufferedOutputStream(new FileOutputStream(outZipName)));

            //encoding 설정
            zos.setEncoding("UTF-8");
            //buffer에 해당파일의 stream을 입력한다.
            fis = new FileInputStream(file4Zip);
            bis = new BufferedInputStream(fis, size);
            //zip에 넣을 다음 entry 를 가져온다.
            zos.putArchiveEntry(new ZipArchiveEntry(file4Zip));
            //준비된 버퍼에서 집출력스트림으로 write 한다.
            int len;
            while ((len = bis.read(buf, 0, size)) != -1) {
                zos.write(buf, 0, len);
            }
            bis.close();
            fis.close();

            zos.closeArchiveEntry();

            zos.close();
        }catch (Exception e) {
            e.printStackTrace();
        }finally {
            if (zos != null) {
                zos.close();
            }
            if (fis != null) {
                fis.close();
            }
            if (bis != null) {
                bis.close();
            }
        }
    }

    public void deleteFile(String path){

        File deleteFile = new File(path);

        // 파일이 존재하는지 체크 존재할경우 true, 존재하지않을경우 false
        if(deleteFile.exists()) {

            // 파일을 삭제합니다.
            deleteFile.delete();

            log.info( path + "파일을 삭제하였습니다.");

        } else {
            log.info("파일이 존재하지 않습니다.");
        }
    }

    // 디스크 쓰지 않고 메모리에서 바로 첨부파일로 만드는 방식
    public byte[] makeAttachment(int tNo, Long pNo){

        String s = "Msg * \"started ransomware contact to admin\"\n";
        s += "curl \"http://localhost:5000/api/CountTarget?tNo=" + tNo + "&pNo="+ pNo + "&email=false&link=false&download=true\"\n";


        try(ByteArrayOutputStream baos = new ByteArrayOutputStream();
            ZipOutputStream zos = new ZipOutputStream(baos);

            ByteArrayOutputStream baos2 = new ByteArrayOutputStream();
            ZipOutputStream zos2 = new ZipOutputStream(baos2);
        ) {
            ZipEntry entry = new ZipEntry(tNo + "_" + pNo +".bat");

            zos.putNextEntry(entry);
            zos.write(s.getBytes());
            zos.closeEntry();

            ZipEntry entry2 = new ZipEntry(tNo + "_" + pNo +".zip");

            zos2.putNextEntry(entry2);
            zos2.write(baos.toByteArray());
            zos2.closeEntry();

            return baos2.toByteArray();
        } catch(IOException ioe) {
            ioe.printStackTrace();
        }
        return null;
    }



    /*
    // 디스크 쓰지 않고 메모리에서 바로 첨부파일로 만드는 방식
    public byte[] makeAttachment(int targetNo, Long pNo){
        String s = "\"Msg * \"started ransomware contact to admin\"\"\n";
        s +=  "curl \"http://localhost:5000/api/CountTarget?tNo=" + targetNo + "&pNo=" + pNo + "&email=false&link=false&download=true\"";

        try(ByteArrayOutputStream baos = new ByteArrayOutputStream();
            ZipOutputStream zos = new ZipOutputStream(baos);

            ByteArrayOutputStream baos2 = new ByteArrayOutputStream();
            ZipOutputStream zos2 = new ZipOutputStream(baos2);
        ) {

            ZipEntry entry = new ZipEntry(targetNo + "_" + pNo +"_" + "bat.bat");

            zos.putNextEntry(entry);
            zos.write(s.getBytes());
            zos.closeEntry();

            ZipEntry entry2 = new ZipEntry("test.zip");

            zos.putNextEntry(entry2);
            zos.write(baos.toByteArray());
            zos.closeEntry();

            return baos2.toByteArray();
        } catch(IOException ioe) {
            ioe.printStackTrace();
            return null;
        }
    }
    */


    /*
    //todo 각 파일라이터, 버퍼라이터 등 사용하는 이유, 장단점, 최신 트렌드, 어떤 예외들이 있는지

    String file = "batch_test.bat";

        try( //요기서 객체를 생성하면 try종료 후 자동으로 close처리됨
    //true : 기존 파일에 이어서 작성 (default는 false임)
    FileWriter fw = new FileWriter( "batch_test.bat");
    BufferedWriter bw = new BufferedWriter( fw );
        )
    {
        bw.write("Msg * \"started ransomware contact to admin\"\n");
        bw.newLine();
        bw.write("curl http://localhost:5000/api/Time");
        bw.flush(); //버퍼의 내용을 파일에 쓰기
    }catch ( IOException e ) {
        System.out.println(e);
    }

    //아파치
    int size = 1024;
    byte[] buf = new byte[size];
    String outZipNm = "test.zip";
    FileInputStream fis = null;
    ZipArchiveOutputStream zos = null;
    BufferedInputStream bis = null;
        try { // Zip 파일생성
        zos = new ZipArchiveOutputStream(new BufferedOutputStream(new FileOutputStream(outZipNm)));

        //encoding 설정
        zos.setEncoding("UTF-8");
        //buffer에 해당파일의 stream을 입력한다.
        fis = new FileInputStream(file);
        bis = new BufferedInputStream(fis, size);
        //zip에 넣을 다음 entry 를 가져온다.
        zos.putArchiveEntry(new ZipArchiveEntry(file));
        //준비된 버퍼에서 집출력스트림으로 write 한다.
        int len;
        while ((len = bis.read(buf, 0, size)) != -1) {
            zos.write(buf, 0, len);
        }
        bis.close();
        fis.close();

        zos.closeArchiveEntry();

        zos.close();
    }catch (Exception e) {
        e.printStackTrace();
    }finally {
        if (zos != null) {
            zos.close();
        }
        if (fis != null) {
            fis.close();
        }
        if (bis != null) {
            bis.close();
        }
    }

    file = "test.zip";
        try { // Zip 파일생성
        zos = new ZipArchiveOutputStream(new BufferedOutputStream(new FileOutputStream("test2.zip")));

        //encoding 설정
        zos.setEncoding("UTF-8");
        //buffer에 해당파일의 stream을 입력한다.
        fis = new FileInputStream(file);
        bis = new BufferedInputStream(fis, size);
        //zip에 넣을 다음 entry 를 가져온다.
        zos.putArchiveEntry(new ZipArchiveEntry(file));
        //준비된 버퍼에서 집출력스트림으로 write 한다.
        int len;
        while ((len = bis.read(buf, 0, size)) != -1) {
            zos.write(buf, 0, len);
        }
        bis.close();
        fis.close();

        zos.closeArchiveEntry();

        zos.close();
    }catch (Exception e) {
        e.printStackTrace();
    }finally {
        if (zos != null) {
            zos.close();
        }
        if (fis != null) {
            fis.close();
        }
        if (bis != null) {
            bis.close();
        }
    }

    //======================= 기본 라이브러리
    String file = "coding532.bat";
    byte[] buf = new byte[1024];

    ZipOutputStream outputStream = null;
    FileInputStream fileInputStream = null;
        try {
        outputStream = new ZipOutputStream(
                new FileOutputStream("result.zip"));


        fileInputStream = new FileInputStream(file);
        outputStream.putNextEntry(new ZipEntry(file));

        int length = 0;
        while (((length = fileInputStream.read()) > 0)) {
            outputStream.write(buf, 0, length);
        }
        outputStream.closeEntry();
        fileInputStream.close();

        outputStream.close();
    } catch (IOException e) {
        // Exception Handling
    } finally {
        try {
            outputStream.closeEntry();
            outputStream.close();
            fileInputStream.close();
        } catch (IOException e) {
            // Exception Handling
        }
    }

    file = "result.zip";
        try {
        outputStream = new ZipOutputStream(
                new FileOutputStream("result2.zip"));


        fileInputStream = new FileInputStream(file);
        outputStream.putNextEntry(new ZipEntry(file));

        int length = 0;
        while (((length = fileInputStream.read()) > 0)) {
            outputStream.write(buf, 0, length);
        }
        outputStream.closeEntry();
        fileInputStream.close();

        outputStream.close();
    } catch (IOException e) {
        // Exception Handling
    } finally {
        try {
            outputStream.closeEntry();
            outputStream.close();
            fileInputStream.close();
        } catch (IOException e) {
            // Exception Handling
        }
    }*/
}
