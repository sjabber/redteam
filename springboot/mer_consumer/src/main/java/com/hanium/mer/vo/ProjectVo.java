package com.hanium.mer.vo;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import javax.persistence.*;
import java.time.LocalDate;
import java.time.LocalDateTime;

@Data
@AllArgsConstructor//
@NoArgsConstructor()//
@Entity(name="project_info")
public class ProjectVo {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "p_no")
    @JsonProperty(value = "p_no")
    private Long pNo;

    @Column(name = "user_no")
    @JsonProperty(value = "user_no")
    private Long userNo;

    @Column(name = "tml_no")
    @JsonProperty(value = "tmp_no")
    private Long tmlNo;


    @Column(name = "tag1")
    @JsonProperty(value = "tag1")
    private int tagFirst;

    @Column(name = "tag2")
    @JsonProperty(value = "tag2")
    private int tagSecond;

    @Column(name = "tag3")
    @JsonProperty(value = "tag3")
    private int tagThird;

    @Column(name = "p_name")
    @JsonProperty(value = "p_name")
    private String pName;

    @Column(name = "p_description")
    @JsonProperty(value = "p_description")
    private String pDescription;

    @Column(name = "p_start_date")
    @JsonProperty(value = "start_date")
    private LocalDate startDate;

    @Column(name = "p_end_date")
    @JsonProperty(value = "end_date")
    private LocalDate endDate;

    @Column(name = "created_time")
    @JsonProperty(value = "created_time")
    private LocalDateTime createdTime;

    @Column(name = "modified_time")
    @JsonProperty(value = "modified_time")
    private LocalDateTime modifiedTime;

    @Column(name = "send_no")
    @JsonProperty(value = "send_no")
    private int sendNo;

    @Column(name = "p_status")
    @JsonProperty(value = "p_status")
    private int pStatus;

    @Column(name = "un_send_no")
    @JsonProperty(value = "un_send_no")
    private int unSendNo;

    @Column(name = "sender_email")
    @JsonProperty(value = "sender_email")
    private String senderEmail;

}
