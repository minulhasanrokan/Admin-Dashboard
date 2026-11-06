'use strict';

(function ($) {
    let designationGrid = null;
    let selectedDesignation = null;
    let currentDesignationUuid = null;

    window.actual_status = {
        1: "New Add",
        2: "Add-Edit",
        3: "Edit",
        4: "Delete",
        5: "Deactivate",
        6: "Active"
    };

    window.process_status = {
        0: "",
        1: "(Sent to Checker)",
        2: "(Returned)",
        3: "(Approved)"
    };

    window.data_status = {
        0: "Deactivate",
        1: "Active"
    };

    window.duty_carry_out = window.duty_carry_out || {};

    if (typeof window.super_admin_status === 'undefined') {
        window.super_admin_status = 1;
    }

    // Configure CSRF for jQuery AJAX
    let csrfToken = $('meta[name="_csrf"]').attr('content') || $('meta[name="csrf-token"]').attr('content') || '';
    if (csrfToken) {
        $.ajaxSetup({
            headers: {
                'X-CSRF-Token': csrfToken
            }
        });
    }

    $(document).ready(function () {
        initDesignationGrid();
        bindFormSubmissions();
        bindGlobalEvents();
    });

    function showNotify(type, message) {
        if (typeof showAlertBox === 'function') {
            let title = type === 'success' ? 'Success' : (type === 'warning' ? 'Warning' : 'Error');
            showAlertBox({
                title: title,
                message: message,
                type: type === 'danger' ? 'danger' : (type === 'warning' ? 'warning' : 'theme')
            });
        } else if (typeof toastNotification === 'function') {
            toastNotification(type, message);
        } else {
            alert(message);
        }
    }

    function getGridRights() {

        let raw = $('#grid_rights_json').val() || '';

        let list = [];

        if (raw) {

            try {
                list = JSON.parse(raw);
            } 
            catch (e) {
                console.warn('Failed to parse grid rights JSON', e);
            }
        }

        return list;
    }

    function createActionColumnDef(right) {

        let action = right.action;
        let actionName = right.action_name;
        let icon = right.icon;
        //let name = right.name;
        //let popUpStatus = right.popup_status;
        let shortName = right.short_name;
        //let sortOrder = right.sort_order;
        let title = right.title;
        let routeUrl = right.url;
        //let width = right.width;

        let headerIconHtml = `<i class="${icon} fs-15"></i>`;

        return {
            headerName: shortName,
            headerIconHtml: headerIconHtml,
            headerIconTitle: title,
            field: 'action_' + actionName,
            colId: 'action_' + actionName,
            width: 38,
            minWidth: 32,
            maxWidth: 44,
            pinned: 'left',
            sortable: false,
            filter: false,
            resizable: false,
            suppressSizeToFit: true,
            headerClass: 'ag-action-col-header text-center',
            cellClass: 'ag-action-col-cell d-flex align-items-center justify-content-center',
            cellRenderer: function (params) {
                let full = params.data;
                if (!full) return '';
                let targetId = full.uuid || full.id;
                let titleValue = encodeURIComponent(full.name);

                let canRender = false;

                if (actionName === 'view') {
                    canRender = true;
                } else if (actionName === 'edit') {
                    if ((full.process_status === 2 || full.process_status === 3 || full.process_status === '0' || full.process_status === '')) {
                        if (full.process_status === 3 || full.process_status === 2 || full.actual_status !== 4) {
                            canRender = true;
                        }
                    }
                } else if (actionName === 'delete') {
                    if ((full.process_status === 2 || full.process_status === 3 || full.process_status === '0' || full.process_status === '')) {
                        if (full.process_status === 3 || full.process_status === 2 || full.actual_status !== 3) {
                            canRender = true;
                        }
                    }
                } else if (actionName === 'stc') {
                    if (full.process_status !== 3 && full.process_status !== 1 && full.actual_status !== 5 && full.actual_status !== 6) {
                        canRender = true;
                    }
                } else if (actionName === 'verify') {
                    let canVerify = (typeof window.duty_carry_out[full.checker] !== 'undefined' || 1 === window.super_admin_status);
                    if (full.process_status === 1 && canVerify) {
                        canRender = true;
                    }
                } else if (actionName === 'active') {
                    if ((full.process_status === 3 || full.process_status === 2 || full.actual_status === 5) && full.status == 0) {
                        canRender = true;
                    }
                } else if (actionName === 'deactivate') {
                    if ((full.process_status === 3 || full.process_status === 2 || full.actual_status === 6) && full.status === 1) {
                        canRender = true;
                    }
                } else {
                    canRender = true;
                }

                if (!canRender) {
                    return '';
                }

                let iconInner = `<i class="${icon} fs-16"></i>`;

                return `<a href="#" class="btn-grid-row-action" title="${title}" data-action="${actionName}" data-url="${routeUrl}" data-title="${title}" data-id="${targetId}" data-name="${titleValue}">${iconInner}</a>`;
            }
        };
    }

    function initDesignationGrid() {
        let rightsList = getGridRights();

        let actionColumnDefs = rightsList.map(createActionColumnDef);

        let dataColumnDefs = [
            {
                field: 'name',
                colId: 'name',
                headerName: 'Designation Name',
                width: 250,
                minWidth: 200,
                sortable: true,
                filter: true,
                cellRenderer: function (params) {
                    return `<span class="fw-semibold text-dark"><i class="ti ti-briefcase me-1 text-primary"></i>${params.data && params.data.name}</span>`;
                }
            },
            {
                field: 'code',
                colId: 'code',
                headerName: 'Code',
                width: 100,
                minWidth: 80,
                sortable: true,
                filter: true,
                cellRenderer: function (params) {
                    let val = params.value || (params.data && params.data.designation_code);
                    return val ? `<span class="badge bg-light text-dark border">${val}</span>` : '<span class="text-muted">-</span>';
                }
            },
            {
                field: 'title',
                colId: 'title',
                headerName: 'Title',
                width: 200,
                minWidth: 160,
                sortable: true,
                filter: true
            },
            {
                field: 'actual_status',
                colId: 'actual_status',
                headerName: 'OP-Status',
                width: 180,
                minWidth: 150,
                sortable: true,
                filter: true,
                cellRenderer: function (params) {
                    let full = params.data || {};
                    let statusText = '';
                    if (typeof window.process_status[full.process_status] !== 'undefined') {
                        statusText = ' ' + window.process_status[full.process_status];
                    }
                    let actLabel = window.actual_status[params.value] || 'New Add';
                    let badgeClass = 'bg-primary-transparent text-primary';
                    if (params.value === 4) badgeClass = 'bg-danger-transparent text-danger';
                    else if (params.value === 5) badgeClass = 'bg-dark-transparent text-dark';
                    else if (params.value === 6) badgeClass = 'bg-success-transparent text-success';
                    else if (params.value === 2) badgeClass = 'bg-warning-transparent text-warning';
                    return `<span class="badge ${badgeClass} fw-semibold">${actLabel}${statusText}</span>`;
                }
            },
            {
                field: 'status',
                colId: 'status',
                headerName: 'Status',
                width: 120,
                minWidth: 100,
                sortable: true,
                filter: true,
                cellRenderer: function (params) {
                    let st = params.value;
                    let badgeClass = st === 1 ? 'bg-success-transparent text-success' : 'bg-danger-transparent text-danger';
                    return `<span class="badge ${badgeClass} fw-semibold">${window.data_status[st] || (st === 1 ? 'Active' : 'Deactivate')}</span>`;
                }
            },
            {
                field: 'message',
                colId: 'message',
                headerName: 'Message',
                width: 220,
                minWidth: 180,
                sortable: true,
                filter: true,
                cellRenderer: function (params) {
                    let full = params.data || {};
                    let msg = params.value;
                    if (!msg) return '<span class="text-muted">-</span>';
                    let safeName = encodeURIComponent(full.name || full.name || '');
                    let safeId = full.uuid || full.id;
                    return `<span class="text-primary cursor-pointer text-truncate d-inline-block mw-100 btn-message-history" data-name="${safeName}" data-id="${safeId}" data-table="designations" title="${msg}">${msg}</span>`;
                }
            },
            {
                field: 'created_at',
                colId: 'created_at',
                headerName: 'Created At',
                width: 180,
                minWidth: 160,
                sortable: true,
                filter: true,
                filterType: 'date',
                cellRenderer: function (params) {
                    return `<span class="text-muted fs-12"><i class="ti ti-clock me-1"></i>${params.value || '-'}</span>`;
                }
            }
        ];

        let allColumnDefs = [
            {
                headerName: '#',
                isIndexCol: true,
                width: 40,
                minWidth: 34,
                maxWidth: 48,
                pinned: 'left',
                sortable: false,
                filter: false,
                resizable: false,
                suppressSizeToFit: true,
                headerClass: 'ag-action-col-header text-center',
                cellClass: 'ag-action-col-cell d-flex align-items-center justify-content-center text-muted fw-medium'
            },
            ...actionColumnDefs,
            ...dataColumnDefs
        ];

        designationGrid = initAgDataGrid({
            gridContainer: '#designation_ag_grid',
            url: '/reference/designation/view',
            exportFileName: 'designations_list',
            defaultOrderCol: 'id',
            defaultOrderDir: 'desc',
            totalBadgeSelector: '#designation_total_badge',
            rowSelection: 'single',
            onSelectionChanged: function (params) {
                if (params && params.api) {
                    let selectedRows = params.api.getSelectedRows();
                    selectedDesignation = selectedRows && selectedRows.length > 0 ? selectedRows[0] : null;
                    if (selectedDesignation) {
                        currentDesignationUuid = selectedDesignation.uuid || selectedDesignation.id;
                    }
                }
            },
            columnDefs: allColumnDefs
        });
    }

    function refreshGrid() {
        if (designationGrid && designationGrid.gridApi) {
            designationGrid.gridApi.refreshServerSide({ purge: true });
        }
    }


    window.return_message_history = function (name, id, table) {
        name = name ? decodeURIComponent(name) : 'Designation';
        openHistoryModal(id, name);
    };



    function openAddModal() {
        openMenu('/reference/designation/add', 'Add Designation', '', 'add', 1, 800, 'popup_1');
    }

    function openEditModal(targetId) {
        currentDesignationUuid = targetId;
        openMenu('/reference/designation/edit', 'Edit Designation', targetId, 'edit', 1, 800, 'popup_1');
    }

    function openViewModal(targetId) {
        currentDesignationUuid = targetId;
        openMenu('/reference/designation/view', 'Designation Details', targetId, 'view', 1, 800, 'popup_1');
    }

    function openDeleteModal(targetId, name) {
        currentDesignationUuid = targetId;
        openMenu('/reference/designation/delete', 'Delete Designation', targetId, 'delete', 1, 600, 'popup_1');
    }

    function openStcModal(targetId, name) {
        currentDesignationUuid = targetId;
        openMenu('/reference/designation/stc', 'Sent To Checker', targetId, 'stc', 1, 800, 'popup_1');
    }

    function openVerifyModal(targetId, name) {
        currentDesignationUuid = targetId;
        openMenu('/reference/designation/verify', 'Verify Designation', targetId, 'verify', 1, 800, 'popup_1');
    }

    function openStatusModal(actionType, targetId, name) {
        currentDesignationUuid = targetId;
        if (actionType === 'active') {
            openMenu('/reference/designation/active', 'Active Designation', targetId, 'active', 1, 600, 'popup_1');
        } else {
            openMenu('/reference/designation/deactivate', 'Deactivate Designation', targetId, 'deactivate', 1, 600, 'popup_1');
        }
    }

    function openHistoryModal(targetId, name) {
        targetId = targetId || currentDesignationUuid;
        if (!targetId) {
            showNotify('warning', 'Please select a designation first');
            return;
        }
        currentDesignationUuid = targetId;
        openMenu('/reference/designation/history', 'Designation Audit History', targetId, 'history', 1, 1140, 'popup_1');
    }

    function openActivityModal(targetId) {
        targetId = targetId || currentDesignationUuid;
        if (!targetId) {
            showNotify('warning', 'Please select a designation first');
            return;
        }
        currentDesignationUuid = targetId;
        openMenu('/reference/designation/activity', 'Designation Activity Logs', targetId, 'activity', 1, 1140, 'popup_1');
    }

    function bindGlobalEvents() {
        // Delegated click handler for AG Grid row actions
        $(document).off('click.designationRow', '.btn-grid-row-action').on('click.designationRow', '.btn-grid-row-action', function (e) {
            e.preventDefault();
            let action = $(this).data('action');
            let id = $(this).data('id');
            let name = $(this).data('name');
            let url = $(this).data('url') || '';

            currentDesignationUuid = id;

            if (!action && url) {
                if (url.indexOf('/view') !== -1) action = 'view';
                else if (url.indexOf('/edit') !== -1) action = 'edit';
                else if (url.indexOf('/delete') !== -1) action = 'delete';
                else if (url.indexOf('/stc') !== -1) action = 'stc';
                else if (url.indexOf('/verify') !== -1) action = 'verify';
                else if (url.indexOf('/active') !== -1) action = 'active';
                else if (url.indexOf('/deactivate') !== -1) action = 'deactivate';
            }

            switch (action) {
                case 'view':
                    openViewModal(id);
                    break;
                case 'edit':
                    openEditModal(id);
                    break;
                case 'delete':
                    openDeleteModal(id, name);
                    break;
                case 'stc':
                    openStcModal(id, name);
                    break;
                case 'verify':
                    openVerifyModal(id, name);
                    break;
                case 'active':
                    openStatusModal('active', id, name);
                    break;
                case 'deactivate':
                    openStatusModal('deactivate', id, name);
                    break;
                default:
                    console.warn('Unhandled action:', action, url);
                    break;
            }
        });

        // Open history from single view modal
        $(document).off('click.designationViewHistory', '#btn_open_history_from_view').on('click.designationViewHistory', '#btn_open_history_from_view', function (e) {
            e.preventDefault();
            let uuid = $(this).data('uuid') || currentDesignationUuid;
            $('#modal_designation_single_view, #modal_designation_view').modal('hide');
            openHistoryModal(uuid, 'Designation History');
        });

        // Activity & History buttons on top right of all sub-modals
        $(document).off('click.designationActivity', '.btn-open-activity').on('click.designationActivity', '.btn-open-activity', function (e) {
            e.preventDefault();
            let $actModal = $('#modal_designation_activity');
            if ($actModal.length && $actModal.find('tbody tr').length > 0) {
                $('.modal.show').modal('hide');
                $actModal.modal('show');
            } else {
                openActivityModal(currentDesignationUuid);
            }
        });

        $(document).off('click.designationHistory', '.btn-open-history').on('click.designationHistory', '.btn-open-history', function (e) {
            e.preventDefault();
            let $histModal = $('#modal_designation_history');
            if ($histModal.length && $histModal.find('tbody tr').length > 0) {
                $('.modal.show').modal('hide');
                $histModal.modal('show');
            } else {
                openHistoryModal(currentDesignationUuid, 'Designation History');
            }
        });

        // Delegated click handler for audit message popup
        $(document).off('click.designationMsg', '.btn-message-history').on('click.designationMsg', '.btn-message-history', function (e) {
            e.preventDefault();
            let name = $(this).data('name');
            let id = $(this).data('id');
            let table = $(this).data('table') || 'designations';
            if (typeof window.return_message_history === 'function') {
                window.return_message_history(name, id, table);
            }
        });
    }

    function validateField(id) {

        switch (id) {
            case "name":
                form_validation({control: "name#text#mandatory", msg_text: "Name"});
                break;
            case "code":
                form_validation({control: "code#text#mandatory", msg_text: "Code"});
                break;
            case "title":
                form_validation({control: "title#text#mandatory", msg_text: "Title"});
                break;
            case "details":
                form_validation({control: "details#text#nonmandatory", msg_text: "Details"});
                break;
            case "delete_reason":
                form_validation({control: "delete_reason#text#mandatory", msg_text: "Delete Reason"});
                break;
            case "checker":
                form_validation({control: "checker#select#mandatory", msg_text: "Checker User"});
                break;
            case "remarks":
                form_validation({control: "remarks#text#nonmandatory", msg_text: "Remarks"});
                break;
            case "return_reason":
                form_validation({control: "return_reason#text#mandatory", msg_text: "Return Reason"});
                break;
            case "active_reason":
                form_validation({control: "active_reason#text#mandatory", msg_text: "Active Reason"});
                break;
            case "deactivate_reason":
                form_validation({control: "deactivate_reason#text#mandatory", msg_text: "Deactivate Reason"});
                break;
        }
    }

    async function saveDesignationAdd() {

        let status = form_validation({control: "name#text#mandatory*code#text#mandatory*title#text#mandatory*details#text#nonmandatory", msg_text: "Name*Code*Title*Details"});

        if (status === false) {

            showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
            return false;
        }

        let confirmed = await confirmDialog("Designation", "Are you sure you want to save the Designation?");

        if (!confirmed) {
            return;
        }

        let formData = new FormData();

        formData.append("name", $("#name").val());
        formData.append("code", $("#code").val());
        formData.append("title", $("#title").val());
        formData.append("details", $("#details").val());

        try {
            let response = await web.post('/reference/designation/add', formData);

            if (response.status === 200) {

                showAlertBox({title: "Success", message: "Saved Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
                $('#modal_designation_add').modal('hide');
                refreshGrid();
            }
        } catch (err) {
        }
    }

    async function saveDesignationEdit() {

        let status = form_validation({control: "name#text#mandatory*code#text#mandatory*title#text#mandatory*details#text#nonmandatory", msg_text: "Name*Code*Title*Details"});

        if (status === false) {

            showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
            return false;
        }

        let confirmed = await confirmDialog("Designation", "Are you sure you want to update the Designation?");

        if (!confirmed) {
            return;
        }

        let uuid = $('#edit_designation_uuid').val() || currentDesignationUuid;

        let formData = new FormData();

        formData.append("name", $("#name").val());
        formData.append("code", $("#code").val());
        formData.append("title", $("#title").val());
        formData.append("details", $("#details").val());

        try {
            let response = await web.post(`/reference/designation/edit/${uuid}`, formData);

            if (response.status === 200) {

                showAlertBox({title: "Success", message: "Updated Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
                $('#modal_designation_edit').modal('hide');
                refreshGrid();
            }
        } catch (err) {
        }
    }

    async function saveDesignationDelete() {

        let status = form_validation({control: "delete_reason#text#mandatory", msg_text: "Delete Reason"});

        if (status === false) {

            showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
            return false;
        }

        let confirmed = await confirmDialog("Delete Designation", "Are you sure you want to delete this Designation?");

        if (!confirmed) {
            return;
        }

        let uuid = $('#delete_target_id').val() || currentDesignationUuid;

        let formData = new FormData();

        formData.append("delete_reason", $("#delete_reason").val());

        try {
            let response = await web.post(`/reference/designation/delete/${uuid}`, formData);

            if (response.status === 200) {

                showAlertBox({title: "Success", message: "Deleted Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
                $('#modal_designation_delete').modal('hide');
                refreshGrid();
            }
        } catch (err) {
        }
    }

    async function saveDesignationStc() {

        let status = form_validation({control: "checker#select#mandatory*remarks#text#nonmandatory", msg_text: "Checker User*Remarks"});

        if (status === false) {

            showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
            return false;
        }

        let confirmed = await confirmDialog("Sent To Checker", "Are you sure you want to send this Designation to Checker?");

        if (!confirmed) {
            return;
        }

        let uuid = $('#stc_target_id').val() || currentDesignationUuid;

        let formData = new FormData();

        formData.append("checker", $("#checker").val());
        formData.append("remarks", $("#remarks").val());

        try {
            let response = await web.post(`/reference/designation/stc/${uuid}`, formData);

            if (response.status === 200) {

                showAlertBox({title: "Success", message: "Sent To Checker Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
                $('#modal_designation_stc').modal('hide');
                refreshGrid();
            }
        } catch (err) {
        }
    }

    async function saveDesignationVerify(verifyType) {

        if (verifyType === 0) {

            let status = form_validation({control: "return_reason#text#mandatory", msg_text: "Return Reason"});

            if (status === false) {

                showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
                return false;
            }

            let confirmed = await confirmDialog("Return Designation", "Are you sure you want to return this Designation?");

            if (!confirmed) {
                return;
            }

            let uuid = $('#verify_target_id').val() || currentDesignationUuid;

            let formData = new FormData();

            formData.append("verify_type", "0");
            formData.append("return_reason", $("#return_reason").val());

            try {
                let response = await web.post(`/reference/designation/verify/${uuid}`, formData);

                if (response.status === 200) {

                    showAlertBox({title: "Success", message: "Returned Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
                    $('#modal_designation_verify').modal('hide');
                    refreshGrid();
                }
            } catch (err) {
            }

        } else {

            let confirmed = await confirmDialog("Verify Designation", "Are you sure you want to approve this Designation?");

            if (!confirmed) {
                return;
            }

            let uuid = $('#verify_target_id').val() || currentDesignationUuid;

            let formData = new FormData();

            formData.append("verify_type", "1");

            try {
                let response = await web.post(`/reference/designation/verify/${uuid}`, formData);

                if (response.status === 200) {

                    showAlertBox({title: "Success", message: "Approved Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
                    $('#modal_designation_verify').modal('hide');
                    refreshGrid();
                }
            } catch (err) {
            }
        }
    }

    async function saveDesignationActive() {

        let status = form_validation({control: "active_reason#text#mandatory", msg_text: "Active Reason"});

        if (status === false) {

            showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
            return false;
        }

        let confirmed = await confirmDialog("Active Designation", "Are you sure you want to activate this Designation?");

        if (!confirmed) {
            return;
        }

        let uuid = $('#active_target_id').val() || currentDesignationUuid;

        let formData = new FormData();

        formData.append("active_reason", $("#active_reason").val());

        try {
            let response = await web.post(`/reference/designation/active/${uuid}`, formData);

            if (response.status === 200) {

                showAlertBox({title: "Success", message: "Activated Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
                $('#modal_designation_active').modal('hide');
                refreshGrid();
            }
        } catch (err) {
        }
    }

    async function saveDesignationDeactivate() {

        let status = form_validation({control: "deactivate_reason#text#mandatory", msg_text: "Deactivate Reason"});

        if (status === false) {

            showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
            return false;
        }

        let confirmed = await confirmDialog("Deactivate Designation", "Are you sure you want to deactivate this Designation?");

        if (!confirmed) {
            return;
        }

        let uuid = $('#deactive_target_id').val() || currentDesignationUuid;

        let formData = new FormData();

        formData.append("deactivate_reason", $("#deactivate_reason").val());

        try {
            let response = await web.post(`/reference/designation/deactivate/${uuid}`, formData);

            if (response.status === 200) {

                showAlertBox({title: "Success", message: "Deactivated Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
                $('#modal_designation_deactive').modal('hide');
                refreshGrid();
            }
        } catch (err) {
        }
    }

    function bindFormSubmissions() {

        $(document).off(".designationForm");

        $(document).on("focusout.designationForm", "input, textarea, select", function () {
            validateField(this.id);
        });

        $(document).on("select2:close.designationForm", "select", function () {
            validateField(this.id);
        });

        $(document).on("click.designationForm", "#btn_save_designation_add", async function () {
            await saveDesignationAdd();
        });

        $(document).on("click.designationForm", "#btn_update_designation", async function () {
            await saveDesignationEdit();
        });

        $(document).on("click.designationForm", "#btn_confirm_delete", async function () {
            await saveDesignationDelete();
        });

        $(document).on("click.designationForm", "#btn_confirm_stc", async function () {
            await saveDesignationStc();
        });

        $(document).on("click.designationForm", "#btn_verify_approve", async function () {
            await saveDesignationVerify(1);
        });

        $(document).on("click.designationForm", "#btn_verify_return", async function () {
            await saveDesignationVerify(0);
        });

        $(document).on("click.designationForm", "#btn_confirm_active", async function () {
            await saveDesignationActive();
        });

        $(document).on("click.designationForm", "#btn_confirm_deactive", async function () {
            await saveDesignationDeactivate();
        });

        $(document).on("submit.designationForm", "#form_designation_add", async function (e) {
            e.preventDefault();
            await saveDesignationAdd();
        });

        $(document).on("submit.designationForm", "#form_designation_edit", async function (e) {
            e.preventDefault();
            await saveDesignationEdit();
        });

        $(document).on("submit.designationForm", "#form_designation_delete", async function (e) {
            e.preventDefault();
            await saveDesignationDelete();
        });

        $(document).on("submit.designationForm", "#form_designation_stc", async function (e) {
            e.preventDefault();
            await saveDesignationStc();
        });

        $(document).on("submit.designationForm", "#form_designation_active", async function (e) {
            e.preventDefault();
            await saveDesignationActive();
        });

        $(document).on("submit.designationForm", "#form_designation_deactive", async function (e) {
            e.preventDefault();
            await saveDesignationDeactivate();
        });
    }

    window.destroyPage = function () {

        $(document).off(".designationForm");
        $(document).off(".designationRow");
        $(document).off(".designationViewHistory");
        $(document).off(".designationActivity");
        $(document).off(".designationHistory");
        $(document).off(".designationMsg");
    };
})(jQuery);
