function Logout() {
    const r = new XMLHttpRequest();
    r.open('GET', 'http://localhost:5000/api/Logout', true);
    // r.setRequestHeader("Content-Type", "application/json");
    r.withCredentials = true;
    r.onreadystatechange = function () {
        // let responseObj;
        if (r.readyState === 4) {
            if (r.status === 200) {
                alert("로그아웃 되었습니다.");
                document.location.href = '/';
            } else {
                console.log("로그아웃 실패");
            }
        }
    };

    r.send();
}

// function CheckLogin() {
//     const r = new XMLHttpRequest();
//     r.open('GET', 'http://localhost:5000/api/CheckLogin', true);
//     r.withCredentials = true;
//     r.onreadystatechange = function () {
//         let responseObj;
//         if (r.readyState === 4) {
//             if (r.status === 200) {
//                 console.log("checkLogin 정상");
//                 responseObj = JSON.parse(r.responseText);
//                 $('#dropdown_id').text(responseObj.user_info.email);
//                 $('#dropdown_name').text(responseObj.user_info.name);
//                 document.location.href = "./dashboard/main.html"
//             } else {
//                 document.location.href = '/';
//             }
//         }
//
//     };
//     r.send();
// }

function CheckLoginInLoginPage() {
    const r = new XMLHttpRequest();
    r.open('GET', 'http://localhost:5000/api/checklogin', true);
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
                document.location.href = '/';
            }
        }

    };
    r.send();
}


function Login() {
    const r = new XMLHttpRequest();
    r.open('POST', 'http://localhost:5000/api/Login', true);
    r.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    r.withCredentials = true;
    r.onreadystatechange = function () {
        if (r.readyState === 4) {
            if (r.status === 200) {
                document.location.href = './dashboard/total'
            } else if (r.status === 400) {
                alert("계정 정보를 입력해주세요. ");
            } else if (r.status === 401) {
                alert("해당 계정이 존재하지 않습니다. ");
            } else if (r.status === 402) {
                alert("패스워드가 일치하지 않습니다. ")
            }
        }
    };
    const formData = $('form').serializeObject();
    r.send(JSON.stringify(formData));
}

function GetDashBoard() {
    const user_id = document.getElementById("dropdown_id");
    const user_name = document.getElementById("dropdown_name");
    const r = new XMLHttpRequest();
    r.open('GET', 'http://localhost:5000/api/dashboard', true);
    r.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    r.withCredentials = true;
    r.onreadystatechange = function () {
        let responseObj;
        if (r.readyState === 4) {
            if (r.status === 200) {
                console.log(r.responseText);
                responseObj = JSON.parse(r.responseText);
                user_id.innerHTML = responseObj.email
                user_name.innerHTML = responseObj.name
                // user_id.text(responseObj.email);
                // user_name.text(responseObj.name);
            } else if (r.status === 403) {
                const tokenResult = Refresh();
                if (tokenResult){
                    console.log("true")
                }

            } else {
                document.location.href = '/';
            }
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
                document.location.href = '/'
            } else if (r.status === 400) {
                alert("정보를 입력해주세요. ");
            } else if (r.status === 401) {
                alert("이미 존재하는 이메일 입니다.");
            } else if (r.status === 402) {
                alert("비밀번호나 이메일 형식이 올바르지 않습니다. ");
            } else if (r.status === 403) {
                alert("비밀번호가 일치하지 않습니다.");
            } else if (r.status === 405) {
                alert("계정을 생성하는 도중 오류가 발생하였습니다. ")
            }

        }
    };
    const formData = $('#register_form').serializeObject();
    r.send(JSON.stringify(formData));
}

function getRoot() {
    const r = new XMLHttpRequest();
    r.open('GET', '/', true);
    r.send()
}

function Refresh() {
    const r = new XMLHttpRequest();
    r.open('GET', 'http://localhost:5000/api/RefreshToken', true);
    r.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    r.withCredentials = true;
    r.onreadystatechange = function () {
        if (r.readyState === 4) {
            if (r.status === 200) {
                console.log(r.responseText);
                return true;
            } else {
                document.location.href = "/";
            }
        }
    };
    const formData = $('form').serializeObject();
    r.send(JSON.stringify(formData));
}
