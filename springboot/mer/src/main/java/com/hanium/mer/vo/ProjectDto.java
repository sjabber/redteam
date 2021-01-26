package com.hanium.mer.vo;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDate;
import java.util.ArrayList;

@Data
@AllArgsConstructor//
@NoArgsConstructor
public class ProjectDto {
    private Long tmp_no;

    private ArrayList<String> tag_no;

    private String p_name;

    private String p_description;

    private LocalDate start_date;

    private LocalDate end_date;

    private int p_status;
    
}
