$(document).ready(function(){
    var username = $(".name");
    var emailcode = $(".k3");
    var userpass = $(".k4");
    $(".k2").click(function(){
        var pname = username.val();
        if(!pname) {
            alert("邮箱错误,请检查");
            return;
        }
        $.post("/signup/email",{email:pname},function(data){
            if(data == 'YES')
            {
                $(".k2").text("发送成功");
            }
        })
    });
    $(".k5").click(function(){
        var pname = username.val();
        var code = emailcode.val();
        var password = userpass.val();
        $.post("/signup/up",{emailname:pname,pwd:password,code:code},function(data){
            if(data == 'YES')
            {
                $(".user *").css("display","none");
                $(".user").css("display","none");
                $(".jump").css("display","block");
                var t=setTimeout(function(){
                    window.location.href='/login';
                },400);
            }
            else
            {
                alert("验证码错误,请检查");
                $(".k2").text("发送验证码");
            }
        })
    });
})