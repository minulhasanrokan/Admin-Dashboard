function form_validation({ control, msg_text }) {

    let controlArr = control.split("*");
    let msgArr = msg_text.split("*");

    let final_status = true;

    for (let i = 0; i < controlArr.length; i++) {

        let error_msg = "";
        let status = true;

        let controlParts = controlArr[i].split("#");

        let originalFieldId = controlParts[0];
        let fieldType = controlParts[1];
        let fieldStatus = controlParts[2];

        let fieldName = msgArr[i] || originalFieldId;

        let el = document.getElementById(originalFieldId);

        if (!el) {

            alert(fieldName + " Not Found!");

            return false;
        }

        if (fieldType === "radio") {

            let radioName = el.getAttribute("name");

            if (!radioName) {

                alert(fieldName + " radio field not found!");

                return false;
            }

            let radioGroup = document.querySelectorAll(`input[type="radio"][name="${CSS.escape(radioName)}"]`);

            let radioErrorDiv = document.querySelector("#error_msg_" + CSS.escape(radioName));

            if (radioErrorDiv) {
                radioErrorDiv.innerText = "";
                radioErrorDiv.style.display = "none";
            }

            radioGroup.forEach(function (radio) {radio.classList.remove(
                "is-valid","is-invalid");});

            let checkedRadio = document.querySelector(`input[type="radio"][name="${CSS.escape(radioName)}"]:checked`);

            if (fieldStatus === "mandatory" && !checkedRadio) {

                error_msg = "Please Select " + fieldName + " Value";
                status = false;
            }

            radioGroup.forEach(function (radio) {

                if (status === false) {

                    radio.classList.add("is-invalid");
                }
            });

            if (status === false) {

                final_status = false;

                showErrorMsg(radioName, error_msg);
            }

            continue;
        }

        let errorDiv = document.querySelector(
            "#error_msg_" + CSS.escape(originalFieldId)
        );

        el.classList.remove("is-valid", "is-invalid");

        let select2Container = el.nextElementSibling;

        if (select2Container && select2Container.classList.contains("select2-container")) {

            let selection = select2Container.querySelector(".select2-selection");

            if (selection) {

                selection.classList.remove("is-valid", "is-invalid", "custom-is-valid", "custom-is-invalid", "form-control");
            }
        }

        if (errorDiv) {
            errorDiv.innerText = "";
            errorDiv.style.display = "none";
        }

        let value = typeof el.value === "string" ? el.value.trim() : "";

        let minLength = el.getAttribute("minlength");
        let maxLength = el.getAttribute("maxlength");

        let inputValueLength = value.length;

        if (fieldType === "checkbox") {

            if (fieldStatus === "mandatory" && !el.checked) {

                error_msg = "Please Check " + fieldName + " Value";
                status = false;
            }
        }

        else if (fieldType === "select") {

            if (fieldStatus === "mandatory" && value === "") {

                error_msg = "Please Select " + fieldName + " Value";
                status = false;
            }
        }
        else if (fieldType === "file") {

            let file = el.files && el.files.length > 0 ? el.files[0] : null;

            if (fieldStatus === "mandatory" && !file) {

                error_msg = "Please Select " + fieldName + " File";
                status = false;
            }

            if (status === true && file) {

                let minFileSizeAttr = el.getAttribute("data-min-size");

                let minFileSize = minFileSizeAttr !== null ? parseFloat(minFileSizeAttr) : 0;

                let maxFileSize = parseFloat(el.getAttribute("data-max-size"));

                if (!isNaN(minFileSize) && file.size < minFileSize * 1024 * 1024) {

                    error_msg = fieldName + " file size must be at least " + minFileSize + " MB";
                    status = false;
                }

                if (status === true && !isNaN(maxFileSize) && file.size > maxFileSize * 1024 * 1024) {

                    error_msg = fieldName + " file size cannot exceed " + maxFileSize + " MB";
                    status = false;
                }

                if (status === true) {

                    let accept = el.getAttribute("accept");

                    if (accept && accept.trim() !== "" && accept.trim() !== "*/*") {

                        let acceptedTypes = accept.split(",").map(function (type) {
                            return type.trim().toLowerCase();
                        }).filter(Boolean);

                        let fileType = (file.type || "").toLowerCase();

                        let fileName = file.name.toLowerCase();

                        let fileExtension = "";

                        if (fileName.includes(".")) {

                            fileExtension = "." + fileName.split(".").pop();
                        }

                        let isValidFileType = acceptedTypes.some(
                            function (type) {

                                if (type.endsWith("/*")) {

                                    let mainType =
                                        type.replace("/*", "/");
                                    return fileType.startsWith(mainType);
                                }
                                if (type.startsWith(".")) {

                                    return (fileExtension === type);
                                }
                                return (fileType === type);
                            }
                        );

                        if (!isValidFileType) {

                            error_msg = fieldName + " has an invalid file type";
                            status = false;
                        }
                    }
                }
            }
        }
        else {

            if (fieldStatus === "mandatory" && value === "") {

                error_msg = "Please Fill Up " + fieldName + " Value";
                status = false;
            }

            if ((fieldType === "text" || fieldType === "email" || fieldType === "mobile" || fieldType === "password" || fieldType === "number") && status === true && value !== "") {

                error_msg = min_max_validation(minLength, maxLength, inputValueLength);

                if (error_msg !== "") {

                    error_msg = fieldName + " " + error_msg;
                    status = false;
                }
            }

            if (fieldType === "email" && status === true && value !== "") {

                let emailFormat = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

                if (!emailFormat.test(value)) {

                    error_msg = fieldName + " Value Will Be E-mail Type";
                    status = false;
                }
            }
            else if (fieldType === "mobile" && status === true && value !== "") {

                let bdMobileFormat = /^(?:\+?88)?01[3-9][0-9]{8}$/;

                if (!bdMobileFormat.test(value)) {

                    error_msg = fieldName + " Value Will Be Bangladeshi Mobile Number";
                    status = false;
                }
            }
            else if (fieldType === "number" && status === true && value !== "") {

                let numberFormat = /^[0-9]+$/;

                if (!numberFormat.test(value)) {

                    error_msg = fieldName + " Value Will Be Number";
                    status = false;

                } else {

                    let minValueAttr = el.getAttribute("data-min-value");
                    let maxValueAttr = el.getAttribute("data-max-value");

                    let numberValue = parseFloat(value);

                    if (minValueAttr !== null && minValueAttr !== "") {

                        let minValue = parseFloat(minValueAttr);

                        if (!isNaN(minValue) && numberValue < minValue) {

                            error_msg = fieldName + " Value Must Be At Least " + minValue;
                            status = false;
                        }
                    }

                    if (status === true && maxValueAttr !== null && maxValueAttr !== "") {

                        let maxValue = parseFloat(maxValueAttr);

                        if (!isNaN(maxValue) && numberValue > maxValue) {

                            error_msg = fieldName + " Value Cannot Be Greater Than " + maxValue;
                            status = false;
                        }
                    }
                }
            }
            else if (fieldType === "password" && status === true && value !== "") {

                error_msg = password_validation(value);

                if (error_msg !== "") {

                    error_msg = fieldName + " " + error_msg;
                    status = false;
                }

                if (originalFieldId === "confirm_password" && status === true) {

                    let passwordEl = document.getElementById("password");

                    let confirmPasswordEl = document.getElementById("confirm_password");

                    if (passwordEl && confirmPasswordEl) {

                        let password = passwordEl.value.trim();

                        let confirmPassword = confirmPasswordEl.value.trim();

                        if (password !== confirmPassword) {

                            error_msg = "Confirm password does not match";
                            status = false;
                        }
                    }
                }
            }
        }
        if (status === false) {

            final_status = false;

            showErrorMsg(originalFieldId, error_msg);

            if (el.classList.contains("select2-hidden-accessible")) {

                let select2Selection = el.nextElementSibling?.querySelector(".select2-selection");

                if (select2Selection) {
                    select2Selection.classList.remove("custom-is-valid");
                    select2Selection.classList.add("custom-is-invalid");
                }
            }
            else {

                el.classList.add("is-invalid");
            }
        }
        else {

            if (el.classList.contains("select2-hidden-accessible")) {

                let select2Selection = el.nextElementSibling?.querySelector(".select2-selection");

                if (select2Selection) {select2Selection.classList.remove("custom-is-invalid");
                    select2Selection.classList.add("custom-is-valid");
                }
            }
            else {
                el.classList.add("is-valid");
            }
        }
    }

    return final_status;
}

function showErrorMsg(field_id, error_msg) {

    let el = document.getElementById(field_id);
    let error_filed = $("#error_msg_" + field_id);

    if (el){

        let select2Container = el.nextElementSibling;

        if (select2Container && select2Container.classList.contains("select2-container")) {

            let selection = select2Container.querySelector(".select2-selection");

            if (selection) {
                selection.classList.remove("custom-is-valid", "custom-is-invalid");
                selection.classList.add("custom-is-invalid");
            }
        }

        el.classList.remove("is-valid");
        el.classList.add("is-invalid");
    }

    if (error_filed.length > 0) {

        error_filed.html(error_msg);
        error_filed.show();
    }
}

function min_max_validation(minLength, maxLength, inputLength) {

    let error_msg = "";

    if (minLength != null && inputLength < minLength) {
        error_msg = "Minimum Length " + minLength + " Characters";
    }

    if (maxLength != null && inputLength > maxLength) {
        error_msg = "Maximum Length " + maxLength + " Characters";
    }

    return error_msg;
}

function password_validation(data) {

    let error_msg = "";
    let value = data.trim();

    if (!/[A-Z]/.test(value))
        error_msg += "Password must contain at least one uppercase letter.<br>";

    if (!/[a-z]/.test(value))
        error_msg += "Password must contain at least one lowercase letter.<br>";

    if (!/[0-9]/.test(value))
        error_msg += "Password must contain at least one digit.<br>";

    if (!/[\W_]/.test(value))
        error_msg += "Password must contain at least one special character.<br>";

    return error_msg;
}


function parseErrors(response) {
    let result = {};
    let errors = response?.errors;

    if (Array.isArray(errors)) {
        errors.forEach(err => {
            if (err.key && err.error) {
                result[err.key] = err.error;
            }
        });
    }

    return result;
}