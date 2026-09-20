if (document.getElementById("eye_btn_password")) {
    document.getElementById("eye_btn_password").addEventListener("click", function () {
        showHideValue('password', 'eye_btn_password');
    });
}

if (document.getElementById("email")) {
    document.getElementById("email").addEventListener("focusout", function () {
        form_validation({control: "email#email#mandatory", msg_text: "Email Address"});
    });
}

if (document.getElementById("mobile")) {
    document.getElementById("mobile").addEventListener("focusout", function () {
        form_validation({control: "mobile#mobile#mandatory", msg_text: "Mobile Number"});
    });
}

if (document.getElementById("password")) {
    document.getElementById("password").addEventListener("focusout", function () {
        form_validation({control: "password#password#mandatory", msg_text: "Password"});
    });
}

if (document.getElementById("remember_me")) {
    document.getElementById("remember_me").addEventListener("click", function () {
        form_validation({control: "remember_me#checkbox#nonmandatory", msg_text: "Remember Me"});
    });
}

document.getElementById("login_button").addEventListener("click", async function () {

    let controls = [];
    let msgs = [];

    if (document.getElementById("mobile")) {
        controls.push("mobile#mobile#mandatory");
        msgs.push("Mobile Number");
    } else if (document.getElementById("email")) {
        controls.push("email#email#mandatory");
        msgs.push("Email Address");
    }

    if (document.getElementById("password")) {
        controls.push("password#password#mandatory");
        msgs.push("Password");
    }

    if (document.getElementById("remember_me")) {
        controls.push("remember_me#checkbox#nonmandatory");
        msgs.push("Remember Me");
    }

    let status = form_validation({
        control: controls.join("*"),
        msg_text: msgs.join("*")
    });

    if (status === false) {

        showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});

        return false;
    }

    showLoader();

    let formData = new FormData();

    if (document.getElementById("mobile")) {
        formData.append("mobile", $("#mobile").val());
    }
    if (document.getElementById("email")) {
        formData.append("email", $("#email").val());
    }
    if (document.getElementById("password")) {
        formData.append("password", $("#password").val());
    }
    if (document.getElementById("remember_me")) {
        formData.append("remember_me", $("#remember_me").is(":checked") ? "1" : "0");
    }

    let response = await web.post('/login', formData);

    if (response.status === 200) {

        let resData = (response.data && response.data.data) ? response.data.data : {};

        if (resData.require_otp) {
            showAlertBox({
                title: "OTP Verification",
                message: response.data.message || "OTP sent successfully. Please verify to proceed.",
                type: "info",
                sound: "/assets/sounds/success.mp3"
            });

            setTimeout(function () {
                window.location.href = BaseUrl + "/user-registration/verify-account/" + resData.uuid;
            }, 2000);
        } else {
            showAlertBox({
                title: "Success",
                message: "User Login Success",
                type: "info",
                sound: "/assets/sounds/success.mp3"
            });

            setTimeout(function () {
                window.location.href = BaseUrl + "/dashboard";
            }, 2000);
        }
    }
});