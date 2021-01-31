package com.hanium.mer.repogitory;


import com.hanium.mer.vo.TargetVo;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface TargetRepository extends JpaRepository<TargetVo, Long> {

    public Optional<TargetVo> findByTargetNo(int target_no);
    //public Optional<TargetVo> findByTagFirstOrTagSecondOrTagThird(int tag_no);
    //TODO 태그에 타겟 없음
    //public int countByTagFirstOrTagSecondOrTagThird(int tag_no);
    @Query(value = "SELECT count(target_no)" +
            "FROM target_info" +
            " WHERE (tag1 != 0 AND tag1 = :tag)" +
            " OR (tag2 != 0 AND tag2 = :tag)" +
            " OR (tag3 != 0 AND tag3 = :tag)", nativeQuery = true)
    Long countByTag(@Param("tag") Integer tag);

}
