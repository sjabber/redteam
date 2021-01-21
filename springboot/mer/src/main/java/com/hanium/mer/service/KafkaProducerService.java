package com.hanium.mer.service;

import com.hanium.mer.vo.KafkaMessage;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.SendResult;
import org.springframework.stereotype.Service;
import org.springframework.util.concurrent.ListenableFuture;
import org.springframework.util.concurrent.ListenableFutureCallback;

@Slf4j
@Service
public class KafkaProducerService {
    private static final String TOPIC = "redteam";
    private final KafkaTemplate<String, KafkaMessage> kafkaTemplate;

    @Autowired
    public SMTPService smtpService;

    @Autowired
    public KafkaProducerService(KafkaTemplate kafkaTemplate) {
        this.kafkaTemplate = kafkaTemplate;
    }

    public void sendMessage(KafkaMessage message) {
        log.info(String.format("Produce message : %s", message));

        //send and forget
        //this.kafkaTemplate.send(TOPIC, message);

        //Async
        ListenableFuture<SendResult<String, KafkaMessage>> future =  this.kafkaTemplate.send(TOPIC, message);
        future.addCallback(new ListenableFutureCallback<SendResult<String, KafkaMessage>>() {
            @Override
            public void onFailure(Throwable ex) {
                log.info("failure producer error {} message : {}", message, ex.getMessage());
            }

            @Override
            public void onSuccess(SendResult<String, KafkaMessage> result) {
                log.info("Success producer message {}", result.getProducerRecord());
                log.info("meta data producer {}", result.getRecordMetadata());
                //log.info("data String producer {}", result.toString());
            }
        });

    }
}
