// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.4.25;
pragma experimental ABIEncoderV2;
import "./Table.sol";

contract Tls{

    event CreateResult(int256 count);
    event InsertResult(int256 count);
    event UpdateResult(int256 count);
    //event RemoveResult(int256 count);

    TableFactory public tableFactory;

    //在构造函数中创建五张表 
    constructor() public{
        tableFactory=TableFactory(0x1001);

        tableFactory.createTable('t_pasture',"A","idA,batch_raw,weight_raw,date_raw,id_pasture");
        tableFactory.createTable('t_factory',"B","idB,batch_pro,product_name,checker,check_time,processer,process_time,composition,id_factory");
        tableFactory.createTable('t_logistics',"C","idC,batch_log,trans_name,id_logistics");
        tableFactory.createTable('t_sales',"D","idD,batch_sale,price,sales_time,id_sales");
        tableFactory.createTable('t_trace',"E","idA,idB,idC,idD");
        tableFactory.createTable('t_ids',"F","idA,idB,idC,idD");

    }


    //通用函数
    //初始化四个id
    function init() public returns (int){
        int cnt=3000;
        Table table=tableFactory.openTable("t_ids");
        Entry entry=table.newEntry();
        entry.set("F","F");
        entry.set("idA",cnt);
        entry.set("idB",cnt);
        entry.set("idC",cnt);
        entry.set("idD",cnt);
        int count=table.insert("F",entry);
        if(count==1){
            return 1;
        }
        else{
            return 0;
        }
    }

    //获取某个id
    function get_id(string idStr) public constant returns (int) {
        Table table=tableFactory.openTable("t_ids");

        Condition condition=table.newCondition();
        condition.EQ("F","F");

        Entries entries=table.select("F",condition);

        int _id = entries.get(0).getInt(idStr);
        return _id;
    }

    //某个id自增
    function set_id(string idStr) public returns (int) {
        int _id = get_id(idStr);
        Table table=tableFactory.openTable("t_ids");

        Condition condition=table.newCondition();
        condition.EQ("F","F");

        Entry entry=table.newEntry();
        entry.set(idStr,_id+1);

        int count=table.update("F",entry,condition);
        if(count==1){
            return _id+1;
        }
        else{
            return -1;
        }
    }

    //初始化四个信息块的关系数据
    function set_trace(int idA)public returns(int){
        int cnt =0;
        Table table2=tableFactory.openTable("t_trace");
        Entry entry2=table2.newEntry();
        entry2.set("E","E");
        entry2.set("idA",idA);
        entry2.set("idB",cnt);
        entry2.set("idC",cnt);
        entry2.set("idD",cnt);
        int count2=table2.insert("E",entry2);
        if(count2 ==1){
            return 1;
        }
        else{
            return 0;
        }
    }

    //用idB 获取 idA
    function getByidB_idA(int _idB)public view returns (int){
        Table table2=tableFactory.openTable("t_trace");
        Condition condition2=table2.newCondition();
        condition2.EQ("idB",_idB);
        Entries entries2=table2.select("E",condition2);//实际上只有一条
        Entry entry2=entries2.get(0);
        int _idA=entry2.getInt("idA");
        return _idA;
    }

    //用idA 获取 id
    function getByidA_id(int _idA,string idStr)public view returns (int){
        Table table=tableFactory.openTable("t_trace");
        Condition condition=table.newCondition();
        condition.EQ("idA",_idA);
        Entries entries=table.select("E",condition);//实际上只有一条
        Entry entry=entries.get(0);
        int _id=entry.getInt(idStr);
        return _id;
    }


    //功能函数

    //牧场模块
    // set_pasture：牧场上传信息的函数
    function set_pasture(string batch_raw,int weight_raw,string date_raw,string id_pasture)public returns(int, int){
        int ret_code=-2;

        Table table=tableFactory.openTable("t_pasture");

        //获取并自增idA
        int idA = set_id("idA");
        if(idA == -1){
            return (-2,-1);
        }

        //新建项
        Entry entry=table.newEntry();
        entry.set("A","A");
        entry.set("idA",idA);
        entry.set("batch_raw",batch_raw);
        entry.set("weight_raw",weight_raw);
        entry.set("date_raw",date_raw);
        entry.set("id_pasture",id_pasture);

        int count=table.insert("A",entry);

        int count2 = set_trace(idA);
        if(count==1&&count2==1){
            // 上传成功
            ret_code=0;
        }else{
            // 上传失败
            ret_code=-1;
        }

        emit InsertResult(count);
        emit InsertResult(count2);

        return (ret_code,idA);
    }

    // getByPasId_pasture：获取某牧场所有历史批次信息
    function getByPasId_pasture(string id_pasture) public view returns(int[] memory,string[] memory,int[] memory,string[] memory){
        Table table=tableFactory.openTable("t_pasture");

        Condition condition=table.newCondition();
        condition.EQ("id_pasture",id_pasture);

        Entries entries=table.select("A",condition);

        int[] memory idA_list=new int[](uint(entries.size()));
        string[] memory batch_raw_list=new string[](uint(entries.size()));
        int256[] memory weight_raw_list=new int256[](uint(entries.size()));
        string[] memory date_raw_list=new string[](uint(entries.size()));

        for(int256 i=0;i<entries.size();++i){
            Entry entry=entries.get(i);

            idA_list[uint(i)]=entry.getInt("idA");
            batch_raw_list[uint(i)]=entry.getString("batch_raw");
            weight_raw_list[uint(i)]=entry.getInt("weight_raw");
            date_raw_list[uint(i)]=entry.getString("date_raw");
        }

        return (idA_list,batch_raw_list,weight_raw_list,date_raw_list);
    }


    //加工厂模块
    // getByFacId_factory1：获取某加工厂历史信息中的成品批次信息（前三个字段）
    function getByFacId_factory1(string id_factory) public view returns(int[] memory,string[] memory,string[] memory){
        Table table=tableFactory.openTable("t_factory");

        Condition condition=table.newCondition();
        condition.EQ("id_factory",id_factory);

        Entries entries=table.select("B",condition);

        int[] memory idB_list=new int[](uint(entries.size()));
        string[] memory batch_pro_list=new string[](uint(entries.size()));
        string[] memory product_name_list=new string[](uint(entries.size()));



        for(int256 i=0;i<entries.size();++i){
            Entry entry=entries.get(i);

            idB_list[uint(i)]=entry.getInt("idB");
            batch_pro_list[uint(i)]=entry.getString("batch_pro");
            product_name_list[uint(i)]=entry.getString("product_name");
        }

        return (idB_list,batch_pro_list,product_name_list);
    }

    // getByFacId_factory2：获取某加工厂历史信息中的成品批次信息（中间两个字段）
    function getByFacId_factory2(string id_factory) public view returns(string[] memory,string[] memory){
        Table table=tableFactory.openTable("t_factory");

        Condition condition=table.newCondition();
        condition.EQ("id_factory",id_factory);

        Entries entries=table.select("B",condition);

        string[] memory checker_list=new string[](uint(entries.size()));
        string[] memory check_time_list=new string[](uint(entries.size()));

        for(int256 i=0;i<entries.size();++i){
            Entry entry=entries.get(i);

            checker_list[uint(i)]=entry.getString("checker");
            check_time_list[uint(i)]=entry.getString("check_time");
        }

        return (checker_list,check_time_list);
    }

    // getByFacId_factory3：获取某加工厂历史信息中的成品批次信息（后三个字段）
    function getByFacId_factory3(string id_factory) public view returns(string[] memory,string[] memory,string[] memory){
        Table table=tableFactory.openTable("t_factory");

        Condition condition=table.newCondition();
        condition.EQ("id_factory",id_factory);

        Entries entries=table.select("B",condition);

        string[] memory processer_list=new string[](uint(entries.size()));
        string[] memory process_time_list=new string[](uint(entries.size()));
        string[] memory composition_list=new string[](uint(entries.size()));

        for(int256 i=0;i<entries.size();++i){
            Entry entry=entries.get(i);

            processer_list[uint(i)]=entry.getString("processer");
            process_time_list[uint(i)]=entry.getString("process_time");
            composition_list[uint(i)]=entry.getString("composition");

        }

        return (processer_list,process_time_list,composition_list);
    }

    // getByFacId_trace1：获取某加工厂在溯源表中的所有记录1(idA[])
    function getByFacId_trace1(string id_factory)public view returns(int[] memory){
        Table table=tableFactory.openTable("t_factory");//B

        // 在t_factory中查找对应id_factory的记录，即某加工厂的记录
        Condition condition=table.newCondition();
        condition.EQ("id_factory",id_factory);
        Entries entries=table.select("B",condition);

        // 定义答案数组
        int[] memory idA_list=new int[](uint(entries.size()));
        // int[] memory idB_list=new int[](uint(entries.size()));

        for(int256 i=0;i<entries.size();++i){
            // 获取该加工厂记录中的第i条
            Entry entry=entries.get(i);

            // 在t_trace中查找上述记录中idB值对应的记录
            int _idB=entry.getInt("idB");

            int _idA = getByidB_idA(_idB);
            idA_list[uint(i)]=_idA;
        }
        return (idA_list);
    }

    // getByFacId_trace12：获取某加工厂在溯源表中的所有记录2(idB[])
    function getByFacId_trace2(string id_factory)public view returns(int[] memory){
        Table table=tableFactory.openTable("t_factory");

        // 在t_factory中查找对应id_factory的记录，即某加工厂的记录
        Condition condition=table.newCondition();
        condition.EQ("id_factory",id_factory);
        Entries entries=table.select("B",condition);

        // 定义答案数组
        // int[] memory idA_list=new int[](uint(entries.size()));
        int[] memory idB_list=new int[](uint(entries.size()));

        for(int256 i=0;i<entries.size();++i){
            // 获取该加工厂记录中的第i条
            Entry entry=entries.get(i);
            idB_list[uint(i)]=entry.getInt("idB");

        }
        return (idB_list);
    }

    // getByIdA_pasture：根据idA数组获取对应牧场的信息
    function getByIdA_pasture(int[] _idA) public view returns(string[] memory,int[] memory,string[] memory,string[] memory){
        Table table=tableFactory.openTable("t_pasture");

        // 定义答案数组
        string[] memory batch_raw_list=new string[](uint(_idA.length));
        int256[] memory weight_raw_list=new int256[](uint(_idA.length));
        string[] memory date_raw_list=new string[](uint(_idA.length));
        string[] memory id_pasture_list=new string[](uint(_idA.length));

        for(uint i=0;i<_idA.length;++i){
            // 获取第i个idA
            int ida=_idA[i];

            // 在表t_pasture中找对应idA值的记录
            Condition condition=table.newCondition();
            condition.EQ("idA",ida);
            Entries entries=table.select("A",condition);//实际上只有一条记录
            // 获取这条记录
            Entry entry=entries.get(0);

            // 将该记录的各字段值存入答案数组中
            batch_raw_list[uint(i)]=entry.getString("batch_raw");
            weight_raw_list[uint(i)]=entry.getInt("weight_raw");
            date_raw_list[uint(i)]=entry.getString("date_raw");
            id_pasture_list[uint(i)]=entry.getString("id_pasture");
        }
        return (batch_raw_list,weight_raw_list,date_raw_list,id_pasture_list);
    }

    // getEmpty_pasture：获取所有暂时未补充加工厂信息的原料奶批次idA数组
    function getEmpty_pasture() public view returns(int[] memory){
        Table table=tableFactory.openTable("t_trace");

        Condition condition=table.newCondition();
        condition.EQ("idB",0);
        Entries entries=table.select("E",condition);

        // 定义答案数组
        int[] memory idA_list=new int[](uint(entries.size()));

        for(int i=0;i<entries.size();++i){
            Entry entry=entries.get(i);

            idA_list[uint(i)]=entry.getInt("idA");
        }
        return idA_list;
    }

    // set_factory1：加工厂上传信息1(前3个字段)
    function set_factory1(int _idA,string batch_pro,string product_name,string composition)public returns(int,int){

        Table table=tableFactory.openTable("t_factory");
        // idB++;

        int idB = set_id("idB");
        if(idB == -1){
            return (-3,-1);
        }

        Entry entry=table.newEntry();
        entry.set("B","B");
        entry.set("idB",idB);
        entry.set("batch_pro",batch_pro);
        entry.set("product_name",product_name);
        entry.set("composition",composition);

        int count=table.insert("B",entry);

        Table table2=tableFactory.openTable("t_trace");
        Entry entry2=table2.newEntry();
        entry2.set("idB",idB);

        Condition condition=table2.newCondition();
        condition.EQ("idA",_idA);

        int count2=table2.update("E",entry2,condition);

        emit InsertResult(count);
        emit InsertResult(count2);
        if(count==1&&count2==1){
            // 上传成功
            return (0,idB);
        }else{
            // 上传失败
            return (-1,-1);
        }

        return (-2,idB);
    }

    // set_factory2：加工厂上传信息2(中间2个字段)
    function set_factory2(int _idA,string checker,string check_time)public returns(int,int){
        int ret_code=-2;

        int _idB = getByidA_id(_idA,"idB");

        Table table=tableFactory.openTable("t_factory");

        Condition condition=table.newCondition();
        condition.EQ("idB",_idB);

        Entry entry=table.newEntry();
        entry.set("checker",checker);
        entry.set("check_time",check_time);

        int count=table.update("B",entry,condition);

        if(count==1){
            // 上传成功
            ret_code=0;
        }else{
            // 上传失败
            ret_code=-1;
        }

        emit InsertResult(count);

        return (ret_code,_idB);
    }

    // set_factory3：加工厂上传信息3(后3个字段)
    function set_factory3(int _idA,string processer,string process_time,string id_factory)public returns(int,int){
        int ret_code=-2;

        int _idB = getByidA_id(_idA,"idB");

        Table table=tableFactory.openTable("t_factory");

        Condition condition=table.newCondition();
        condition.EQ("idB",_idB);

        Entry entry=table.newEntry();
        entry.set("processer",processer);
        entry.set("process_time",process_time);
        entry.set("id_factory",id_factory);

        int count=table.update("B",entry,condition);

        if(count==1){
            // 上传成功
            ret_code=0;
        }else{
            // 上传失败
            ret_code=-1;
        }

        emit InsertResult(count);

        return (ret_code,_idB);
    }

    //储运商
    // getByLogid_logistics：获取某储运商历史信息中的运输批次信息
    function getByLogid_logistics(string id_logistics) public view returns(int[] memory,string[] memory,string[] memory,string[] memory){
        Table table=tableFactory.openTable("t_logistics");

        Condition condition=table.newCondition();
        condition.EQ("id_logistics",id_logistics);

        Entries entries=table.select("C",condition);

        int[] memory idC_list=new int[](uint(entries.size()));
        string[] memory trans_name_list=new string[](uint(entries.size()));
        string[] memory batch_log_list=new string[](uint(entries.size()));
        string[] memory id_logistics_list=new string[](uint(entries.size()));

        for(int256 i=0;i<entries.size();++i){
            Entry entry=entries.get(i);

            idC_list[uint(i)]=entry.getInt("idC");
            id_logistics_list[uint(i)]=entry.getString("id_logistics");
            trans_name_list[uint(i)]=entry.getString("trans_name");
            batch_log_list[uint(i)]=entry.getString("batch_log");

        }

        return (idC_list,id_logistics_list,trans_name_list,batch_log_list);
    }

    // getByLogId_trace1：获取某储运商在溯源表中的所有记录1(idA[])
    function getByLogId_trace1(string id_logistics)public view returns(int[] memory){
        Table table=tableFactory.openTable("t_logistics");
        Table table2=tableFactory.openTable("t_trace");

        Condition condition=table.newCondition();
        condition.EQ("id_logistics",id_logistics);
        Entries entries=table.select("C",condition);

        int[] memory idA_list=new int[](uint(entries.size()));

        for(int256 i=0;i<entries.size();++i){
            Entry entry=entries.get(i);

            int _idC=entry.getInt("idC");

            Condition condition2=table2.newCondition();
            condition2.EQ("idC",_idC);
            Entries entries2=table2.select("E",condition2);

            Entry entry2=entries2.get(0);

            idA_list[uint(i)]=entry2.getInt("idA");
        }

        return (idA_list);
    }

    // getByLogId_trace2：获取某储运商在溯源表中的所有记录2(idB[])
    function getByLogId_trace2(string id_logistics)public view returns(int[] memory){
        Table table=tableFactory.openTable("t_logistics");
        Table table2=tableFactory.openTable("t_trace");

        Condition condition=table.newCondition();
        condition.EQ("id_logistics",id_logistics);
        Entries entries=table.select("C",condition);

        int[] memory idB_list=new int[](uint(entries.size()));

        for(int256 i=0;i<entries.size();++i){
            Entry entry=entries.get(i);

            int _idC=entry.getInt("idC");

            Condition condition2=table2.newCondition();
            condition2.EQ("idC",_idC);
            Entries entries2=table2.select("E",condition2);

            Entry entry2=entries2.get(0);

            idB_list[uint(i)]=entry2.getInt("idB");
        }

        return (idB_list);
    }

    // getByLogId_trace3：获取某储运商在溯源表中的所有记录3(idC[])
    function getByLogId_trace3(string id_logistics)public view returns(int[] memory){
        Table table=tableFactory.openTable("t_logistics");
        Table table2=tableFactory.openTable("t_trace");

        Condition condition=table.newCondition();
        condition.EQ("id_logistics",id_logistics);
        Entries entries=table.select("C",condition);

        int[] memory idC_list=new int[](uint(entries.size()));

        for(int256 i=0;i<entries.size();++i){
            Entry entry=entries.get(i);

            int _idC=entry.getInt("idC");

            Condition condition2=table2.newCondition();
            condition2.EQ("idC",_idC);
            Entries entries2=table2.select("E",condition2);

            Entry entry2=entries2.get(0);

            idC_list[uint(i)]=entry2.getInt("idC");
        }

        return (idC_list);
    }

    // getByIdB_factory1：根据idB数组获取对应加工厂的信息(前三个字段)
    function getByIdB_factory1(int[] _idB)public view returns(int[] memory,string[] memory,string[] memory){
        Table table=tableFactory.openTable("t_factory");

        int[] memory idB_list=new int[](uint(_idB.length));
        string[] memory batch_pro_list=new string[](uint(_idB.length));
        string[] memory product_name_list=new string[](uint(_idB.length));

        for(uint i=0;i<_idB.length;++i){
            // 获取第i个idB
            int B=_idB[i];

            // 在表t_factory中找对应idB值的记录
            Condition condition=table.newCondition();
            condition.EQ("idB",B);
            Entries entries=table.select("B",condition);//实际上只有一条记录
            // 获取这条记录
            Entry entry=entries.get(0);

            // 将该记录的各字段值存入答案数组中
            idB_list[uint(i)]=entry.getInt("idB");
            batch_pro_list[uint(i)]=entry.getString("batch_pro");
            product_name_list[uint(i)]=entry.getString("product_name");
        }

        return (idB_list,batch_pro_list,product_name_list);
    }

    //  getByIdB_factory2：根据idB数组获取对应加工厂的信息（中间两个字段）
    function getByIdB_factory2(int[] _idB) public view returns(string[] memory,string[] memory){
        Table table=tableFactory.openTable("t_factory");

        string[] memory checker_list=new string[](uint(_idB.length));
        string[] memory check_time_list=new string[](uint(_idB.length));

        for(uint i=0;i<_idB.length;++i){
            // 获取第i个idB
            int B=_idB[i];

            // 在表t_factory中找对应idB值的记录
            Condition condition=table.newCondition();
            condition.EQ("idB",B);
            Entries entries=table.select("B",condition);//实际上只有一条记录
            // 获取这条记录
            Entry entry=entries.get(0);

            // 将该记录的各字段值存入答案数组中
            checker_list[uint(i)]=entry.getString("checker");
            check_time_list[uint(i)]=entry.getString("check_time");
        }

        return (checker_list,check_time_list);
    }

    // getByIdB_factory3：根据idB数组获取对应加工厂的信息（后三个字段）
    function getByIdB_factory3(int[] _idB)public view returns(string[] memory,string[] memory,string[] memory){
        Table table=tableFactory.openTable("t_factory");

        string[] memory processer_list=new string[](uint(_idB.length));
        string[] memory process_time_list=new string[](uint(_idB.length));
        string[] memory composition_list=new string[](uint(_idB.length));

        for(uint i=0;i<_idB.length;++i){
            // 获取第i个idB
            int B=_idB[i];

            // 在表t_factory中找对应idB值的记录
            Condition condition=table.newCondition();
            condition.EQ("idB",B);
            Entries entries=table.select("B",condition);//实际上只有一条记录
            // 获取这条记录
            Entry entry=entries.get(0);

            // 将该记录的各字段值存入答案数组中
            processer_list[uint(i)]=entry.getString("processer");
            process_time_list[uint(i)]=entry.getString("process_time");
            composition_list[uint(i)]=entry.getString("composition");
        }

        return (processer_list,process_time_list,composition_list);
    }

    // getEmpty_factory：获取所有暂时未补充储运商信息的成品批次的idA和idB数组
    function getEmpty_factory() public view returns(int[] memory,int[] memory){
        Table table=tableFactory.openTable("t_trace");

        Condition condition=table.newCondition();
        condition.EQ("idC",0);
        condition.NE("idB",0);
        Entries entries=table.select("E",condition);

        // 定义答案数组
        int[] memory idA_list=new int[](uint(entries.size()));
        int[] memory idB_list=new int[](uint(entries.size()));

        for(int i=0;i<entries.size();++i){
            Entry entry=entries.get(i);
            idA_list[uint(i)]=entry.getInt("idA");
            idB_list[uint(i)]=entry.getInt("idB");

        }
        return (idA_list,idB_list);
    }

    // set_logistics：储运商上传信息
    function set_logistics(int _idA,string batch_log,string trans_name,string id_logistics)public returns(int,int){

        Table table=tableFactory.openTable("t_logistics");
        // idC++;
        int idC = set_id("idC");
        if(idC == -1){
            return (-3,-1);
        }

        Entry entry=table.newEntry();
        entry.set("idC",idC);
        entry.set("batch_log",batch_log);
        entry.set("trans_name",trans_name);
        entry.set("id_logistics",id_logistics);

        int count=table.insert("C",entry);

        Table table2=tableFactory.openTable("t_trace");
        Entry entry2=table2.newEntry();
        entry2.set("idC",idC);

        Condition condition=table2.newCondition();
        condition.EQ("idA",_idA);

        int count2=table2.update("E",entry2,condition);

        emit InsertResult(count);
        emit InsertResult(count2);

        if(count==1&&count2==1){
            // 上传成功
            return(0,idC);
        }else{
            // 上传失败
            return (-1,-1);
        }
    }

    //销售商
    // getBySaleid_sales：获取某销售商历史信息中的销售批次信息
    function getBySaleid_sales(string id_sales)public view returns(int[] memory,string[] memory,string[] memory,string[] memory){
        Table table=tableFactory.openTable("t_sales");

        Condition condition=table.newCondition();
        condition.EQ("id_sales",id_sales);

        Entries entries=table.select("D",condition);

        int[] memory idD_list=new int[](uint(entries.size()));
        string[] memory price_list=new string[](uint(entries.size()));
        string[] memory sales_time_list=new string[](uint(entries.size()));
        string[] memory batch_sale_list=new string[](uint(entries.size()));

        for(int256 i=0;i<entries.size();++i){
            Entry entry=entries.get(i);

            idD_list[uint(i)]=entry.getInt("idD");
            price_list[uint(i)]=entry.getString("price");
            sales_time_list[uint(i)]=entry.getString("sales_time");
            batch_sale_list[uint(i)]=entry.getString("batch_sale");

        }

        return (idD_list,price_list,sales_time_list,batch_sale_list);
    }

    // getBySaleId_trace1：获取某销售商在溯源表中的所有记录1(idA[])
    function getBySaleId_trace1(string id_sales)public view returns(int[] memory){
        Table table=tableFactory.openTable("t_sales");
        Table table2=tableFactory.openTable("t_trace");

        Condition condition=table.newCondition();
        condition.EQ("id_sales",id_sales);
        Entries entries=table.select("D",condition);

        int[] memory idA_list=new int[](uint(entries.size()));

        for(int i=0;i<entries.size();++i){
            Entry entry=entries.get(i);

            int _idD=entry.getInt("idD");

            Condition condition2=table2.newCondition();
            condition2.EQ("idD",_idD);
            Entries entries2=table2.select("E",condition2);

            Entry entry2=entries2.get(0);

            idA_list[uint(i)]=entry2.getInt("idA");
        }

        return (idA_list);
    }

    // getBySaleId_trace2：获取某销售商在溯源表中的所有记录2(idB[])
    function getBySaleId_trace2(string id_sales)public view returns(int[] memory){
        Table table=tableFactory.openTable("t_sales");
        Table table2=tableFactory.openTable("t_trace");

        Condition condition=table.newCondition();
        condition.EQ("id_sales",id_sales);
        Entries entries=table.select("D",condition);

        int[] memory idB_list=new int[](uint(entries.size()));

        for(int i=0;i<entries.size();++i){
            Entry entry=entries.get(i);

            int _idD=entry.getInt("idD");

            Condition condition2=table2.newCondition();
            condition2.EQ("idD",_idD);
            Entries entries2=table2.select("E",condition2);

            Entry entry2=entries2.get(0);

            idB_list[uint(i)]=entry2.getInt("idB");
        }

        return (idB_list);
    }

    // getBySaleId_trace3：获取某销售商在溯源表中的所有记录3(idC[])
    function getBySaleId_trace3(string id_sales)public view returns(int[] memory){
        Table table=tableFactory.openTable("t_sales");
        Table table2=tableFactory.openTable("t_trace");

        Condition condition=table.newCondition();
        condition.EQ("id_sales",id_sales);
        Entries entries=table.select("D",condition);

        int[] memory idC_list=new int[](uint(entries.size()));

        for(int i=0;i<entries.size();++i){
            Entry entry=entries.get(i);

            int _idD=entry.getInt("idD");

            Condition condition2=table2.newCondition();
            condition2.EQ("idD",_idD);
            Entries entries2=table2.select("E",condition2);

            Entry entry2=entries2.get(0);

            idC_list[uint(i)]=entry2.getInt("idC");
        }

        return (idC_list);
    }

    // getBySaleId_trace4：获取某销售商在溯源表中的所有记录4(idD[])
    function getBySaleId_trace4(string id_sales)public view returns(int[] memory){
        Table table=tableFactory.openTable("t_sales");
        Table table2=tableFactory.openTable("t_trace");

        Condition condition=table.newCondition();
        condition.EQ("id_sales",id_sales);
        Entries entries=table.select("D",condition);

        int[] memory idD_list=new int[](uint(entries.size()));

        for(int i=0;i<entries.size();++i){
            Entry entry=entries.get(i);

            int _idD=entry.getInt("idD");

            Condition condition2=table2.newCondition();
            condition2.EQ("idD",_idD);
            Entries entries2=table2.select("E",condition2);

            Entry entry2=entries2.get(0);

            idD_list[uint(i)]=entry2.getInt("idD");
        }

        return (idD_list);
    }

    // getByIdC_logistics：根据idC数组获取对应储运商的信息
    function getByIdC_logistics(int[] idC_list)public view returns(int[] memory,string[] memory,string[] memory,string[] memory){
        Table table=tableFactory.openTable("t_logistics");

        // 定义答案数组
        string[] memory batch_log_list=new string[](uint(idC_list.length));
        string[] memory trans_name_list=new string[](uint(idC_list.length));
        string[] memory id_logistics_list=new string[](uint(idC_list.length));

        for(uint i=0;i<idC_list.length;++i){
            // 在t_trace中查找上述记录中idD值对应的记录
            Condition condition=table.newCondition();
            condition.EQ("idC",idC_list[i]);
            Entries entries=table.select("C",condition);//实际上只有一条记录
            // 获取这条记录
            Entry entry=entries.get(0);

            // 将该记录的各字段值存入答案数组中
            batch_log_list[uint(i)]=entry.getString("batch_log");
            trans_name_list[uint(i)]=entry.getString("trans_name");
            id_logistics_list[uint(i)]=entry.getString("id_logistics");
        }

        return (idC_list,batch_log_list,trans_name_list,id_logistics_list);

    }

    // getEmpty_sales：获取溯源表中所有暂时未补充销售商信息的记录(idA[],idB[],idC[],idD[]="")
    function getEmpty_logistics()public view returns(int[] memory,int[] memory,int[] memory){
        Table table=tableFactory.openTable("t_trace");

        Condition condition=table.newCondition();
        condition.EQ("idD",0);
        condition.NE("idC",0);
        condition.NE("idB",0);

        Entries entries=table.select("E",condition);

        // 定义答案数组
        int[] memory idA_list=new int[](uint(entries.size()));
        int[] memory idB_list=new int[](uint(entries.size()));
        int[] memory idC_list=new int[](uint(entries.size()));

        for(int i=0;i<entries.size();++i){
            Entry entry=entries.get(i);

            idA_list[uint(i)]=entry.getInt("idA");
            idB_list[uint(i)]=entry.getInt("idB");
            idC_list[uint(i)]=entry.getInt("idC");
        }

        return (idA_list,idB_list,idC_list);
    }

    // set_sales：销售商上传信息的函数
    function set_sales(int _idA,string batch_sale,int price,string sales_time,string id_sales)public returns(int,int){
        int ret_code=0;

        Table table=tableFactory.openTable("t_sales");

        int idD = set_id("idD");
        if(idD == -1){
            return (-3,-1);
        }

        Entry entry=table.newEntry();
        entry.set("D","D");
        entry.set("idD",idD);
        entry.set("batch_sale",batch_sale);
        entry.set("price",price);
        entry.set("sales_time",sales_time);
        entry.set("id_sales",id_sales);

        int count=table.insert("D",entry);

        int count2 = setIdDTrace(_idA,idD);
        if(count2!=0){
            return (-2,-1);
        }

        if(count==1){
            // 上传成功
            ret_code=0;
        }else{
            // 上传失败
            ret_code=-1;
        }

        emit InsertResult(count);
        // emit InsertResult(count2);

        return (ret_code,idD);
    }

    //为trace表增加idD的值
    function setIdDTrace(int _idA,int _idD) public returns(int){
        Table table=tableFactory.openTable("t_trace");
        Entry entry=table.newEntry();
        entry.set("idD",_idD);

        Condition condition=table.newCondition();
        condition.EQ("idA",_idA);
        int count=table.update("E",entry,condition);
        if(count==1){
            // 上传成功
            return 0;
        }else{
            // 上传失败
            return -1;
        }
    }

    //溯源用
    //通过idA获取pasture信息
    function getPasture(int _idA) public view returns(int, string ,int ,string ,string ){
        Table table=tableFactory.openTable("t_pasture");
        Condition condition=table.newCondition();
        condition.EQ("idA",_idA);

        Entries entries=table.select("A",condition);
        Entry entry = entries.get(0);
        return (entry.getInt("idA"),entry.getString("batch_raw"),entry.getInt("weight_raw"),
        entry.getString("date_raw"),entry.getString("id_pasture"));
    }

    //通过idB获取factory信息
    function getFactory(int _idB) public view returns(int,string,string,string,string,string,string,string,string){
        Table table=tableFactory.openTable("t_factory");
        Condition condition=table.newCondition();
        condition.EQ("idB",_idB);

        Entries entries=table.select("B",condition);
        Entry entry = entries.get(0);
        return (entry.getInt("idB"),entry.getString("batch_pro"),entry.getString("product_name"),
        entry.getString("composition"),entry.getString("checker"),entry.getString("check_time"),
        entry.getString("processer"),entry.getString("process_time"),entry.getString("id_factory"));
    }

    //通过idD获取logistics信息
    function getLogistics(int _idC) public view returns(int, string ,string ,string ){
        Table table=tableFactory.openTable("t_logistics");
        Condition condition=table.newCondition();
        condition.EQ("idC",_idC);

        Entries entries=table.select("C",condition);
        Entry entry = entries.get(0);
        return (entry.getInt("idC"),entry.getString("batch_log"),
        entry.getString("trans_name"),entry.getString("id_logistics"));
    }

    //通过idD获取sale信息
    function getSales(int _idD) public view returns(int, string ,string ,string ,string ){
        Table table=tableFactory.openTable("t_sales");
        Condition condition=table.newCondition();
        condition.EQ("idD",_idD);

        Entries entries=table.select("D",condition);
        Entry entry = entries.get(0);
        return (entry.getInt("idD"),entry.getString("batch_sale"),
        entry.getString("price"),entry.getString("sales_time"),entry.getString("id_sales"));
    }

}