const rootUrl = 'http://localhost:5000';


function Logout() {
    const r = new XMLHttpRequest();
    r.open('GET', rootUrl + '/logout', true);
    // r.setRequestHeader("Content-Type", "application/json");
    r.withCredentials = true;
    r.onreadystatechange = function () {
        let responseObj;
        if (r.readyState === 4) {
            if (r.status === 200) {
                document.location.href = "../index.html";
            } else {
                console.log("로그아웃 실패")
            }
        }
    };

    r.send();
}

function CheckLogin() {
    const r = new XMLHttpRequest();
    r.open('GET', 'http://localhost:5000/api/checklogin', true);
    r.withCredentials = true;
    r.onreadystatechange = function () {
        let responseObj;
        if (r.readyState === 4) {
            if (r.status === 200) {
                console.log("checkLogin 정상");
                responseObj = JSON.parse(r.responseText);
                $('#dropdown_id').text(responseObj.user_info.email);
                $('#dropdown_name').text(responseObj.user_info.name);
                // document.location.href = "./dashboard/main.html"
            } else {
                document.location.href = "../index.html";
            }
        }

    };
    r.send();
}

function CheckLoginInLoginPage() {
    const r = new XMLHttpRequest();
    r.open('GET', 'http://localhost:5000/checklogin', true);
    r.withCredentials = true;
    r.onreadystatechange = function () {
        // let responseObj;
        if (r.readyState === 4) {
            if (r.status === 200) {
                console.log("checkLogin 정상");
                // responseObj = JSON.parse(r.responseText);
                // $('#dropdown_id').text(responseObj.user_info.email);
                // $('#dropdown_name').text(responseObj.user_info.name);
                document.location.href = "./dashboard/main.html"
            } else {
                document.location.href = "../index.html";
            }
        }

    };
    r.send();
}


function Login() {
    const r = new XMLHttpRequest();
    r.open('POST', 'http://localhost:5000/api/login', true);
    r.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    r.withCredentials = true;
    r.onreadystatechange = function () {
        if (r.readyState === 4) {
            if (r.status === 200) {
                document.location.href = "./dashboard/main.html"
            } else if (r.status === 400) {
                alert("계정 정보를 입력해주세요");
            } else if (r.status === 401) {
                alert("인증 실패했습니다.");
            }
        }
    };
    const formData = $('form').serializeObject();
    r.send(JSON.stringify(formData));
}

function GetDashBoard() {
    const user_id = document.getElementById("user_id");
    const user_name = document.getElementById("user_name");
    const r = new XMLHttpRequest();
    r.open('POST', 'http://localhost:5000/api/dashboard', true);
    r.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    r.withCredentials = true;
    r.onreadystatechange = function () {
        let responseObj;
        if (r.readyState === 4) {
            if (r.status === 200) {
                // console.log(r.responseText);
                responseObj = JSON.parse(r.responseText);
                $('#user_id').text(responseObj.email);
                $('#user_name').text(responseObj.name);
            } else {
                document.location.href = "../index.html";
            }/checklogin
        }
    };

    r.send();
}

function Register() {
    if (!$("input:checkbox[name='checkbox']").is(":checked")) {
        alert("이용약관에 동의해주세요");
        return;
    }
    const r = new XMLHttpRequest();
    r.open('POST', 'http://localhost:5000/api/createUser', true);
    r.onreadystatechange = function () {
        if (r.readyState === 4) {
            if (r.status === 200) {
                alert("회원가입이 완료되었습니다.");
                document.location.href = "index.html"
            } else if (r.status === 400) {
                alert("정보를 입력해주세요");
            } else if (r.status === 403) {
                alert("이미 존재하는 이메일 입니다.");
            } else if (r.status === 401) {
                alert("비밀번호가 틀렸습니다.");
            } else if (r.status === 402) {
                alert("비밀번호는 8자리 이상 입력해주세요.");
            }
        }
    };
    const formData = $('#register_form').serializeObject();
    r.send(JSON.stringify(formData));
}
