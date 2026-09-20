document.getElementById("forgot_password").addEventListener("click", async function () {

    let status = form_validation({control: "email#email#mandatory", msg_text: "Email Address"});

    if (status === false) {

        showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});

        return false;
    }

    showLoader();

    let formData = new FormData();

    formData.append("email", $("#email").val());

    let response = await web.post('/user-password/forgot-password-store', formData);

    if (response.status === 200) {

        showAlertBox({
            title: "Success",
            message: "Password reset email has been sent successfully",
            type: "info",
            sound: "/assets/sounds/success.mp3"
        });

        setTimeout(function () {
            window.location.href = BaseUrl + "/user-password/forgot-password/" + response.data.data.uuid;
        }, 300);
    }
});

document.getElementById("email").addEventListener("focusout", function () {

    form_validation({control: "email#email#mandatory", msg_text: "Email Address"});
});