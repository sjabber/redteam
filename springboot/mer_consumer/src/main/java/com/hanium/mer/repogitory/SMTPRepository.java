package com.hanium.mer.repogitory;

import com.hanium.mer.vo.SmtpVo;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface SMTPRepository extends JpaRepository<SmtpVo, Long> {

    public Optional<SmtpVo> findByUserNo(Long user_no);

}
