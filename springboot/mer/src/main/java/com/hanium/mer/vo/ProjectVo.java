package com.hanium.mer.vo;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AccessLevel;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import javax.persistence.*;
import java.time.LocalDate;
import java.time.LocalDateTime;

@Data
@AllArgsConstructor//
@NoArgsConstructor(access = AccessLevel.PROTECTED)//
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

    @Column(name = "tmp_no")
    @JsonProperty(value = "tmp_no")
    private Long tmpNo;

    @Column(name = "tag_1")
    @JsonProperty(value = "tag_1")
    private Long tagFirst;

    @Column(name = "tag_2")
    @JsonProperty(value = "tag_2")
    private Long tagSecond;

    @Column(name = "tag_3")
    @JsonProperty(value = "tag_3")
    private Long tagThird;

    @Column(name = "p_name")
    @JsonProperty(value = "p_name")
    private String pName;

    @Column(name = "p_description")
    @JsonProperty(value = "p_description")
    private String pDescription;

    @Column(name = "start_date")
    @JsonProperty(value = "start_date")
    private LocalDate startDate;

    @Column(name = "end_date")
    @JsonProperty(value = "end_date")
    private LocalDate endDate;

    @Column(name = "created_time")
    @JsonProperty(value = "created_time")
    private LocalDateTime createdTime;

    @Column(name = "modified_time")
    @JsonProperty(value = "modified_time")
    private LocalDateTime modifiedTime;

}
