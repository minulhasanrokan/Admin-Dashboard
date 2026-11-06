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

    function closeModal(selector) {
        let el = typeof selector === 'string' ? document.querySelector(selector) : selector;
        if (el) {
            let inst = (typeof bootstrap !== 'undefined' && bootstrap.Modal) ? bootstrap.Modal.getInstance(el) : null;
            if (inst) {
                inst.hide();
            } else if (typeof $(el).modal === 'function') {
                $(el).modal('hide');
            }
        }
    }

    function handleAjaxError(err, defaultMsg) {
        let msg = defaultMsg || "An unexpected error occurred";
        if (err && err.response && err.response.data) {
            let data = err.response.data;
            if (data.error && typeof data.error === 'string' && data.error.length > 0) {
                msg = data.error.replace(/^\[ERROR\]\s*(Internal Server Error:\s*)?/i, '');
            } else if (data.message && data.message !== "Internal Server Error") {
                msg = data.message;
            } else if (data.message) {
                msg = data.message;
            }
            if (Array.isArray(data.errors) && data.errors.length > 0) {
                let errorItems = data.errors.map(function (e) {
                    return e.message || e.field || e;
                }).join(', ');
                msg += ": " + errorItems;
            }
        } else if (err && err.message) {
            msg = err.message;
        }
        showAlertBox({title: "Error", message: msg, type: "danger", sound: "/assets/sounds/error.mp3"});
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

    function createActionsColumnDef(rightsList) {
        let rightsMap = {};
        if (Array.isArray(rightsList)) {
            rightsList.forEach(function (r) {
                if (r && r.action_name) {
                    rightsMap[r.action_name] = r;
                }
            });
        }

        return {
            headerName: 'Actions',
            colId: 'actions',
            field: 'actions',
            width: 175,
            minWidth: 155,
            pinned: 'left',
            sortable: false,
            filter: false,
            resizable: false,
            suppressSizeToFit: true,
            headerClass: 'text-center',
            cellClass: 'ag-action-cell',
            cellRenderer: function (params) {
                let full = params.data;
                if (!full) return '';
                let targetId = full.uuid || full.id;
                let titleValue = encodeURIComponent(full.name || full.designation_name || '');

                let procStatus = Number(full.process_status);
                let actStatus = Number(full.actual_status);
                let recStatus = Number(full.status);

                let buttons = [];

                // View
                if (rightsMap['view']) {
                    let r = rightsMap['view'];
                    let icon = r.icon || 'ti ti-eye';
                    buttons.push(`<a href="#" class="btn-grid-row-action text-info" title="${r.title || 'View Designation'}" data-action="view" data-url="${r.url || ''}" data-title="${r.title || 'View Designation'}" data-id="${targetId}" data-name="${titleValue}"><i class="${icon} fs-16"></i></a>`);
                }

                // Edit
                if (rightsMap['edit'] && (procStatus === 0 || procStatus === 2 || procStatus === 3) && actStatus !== 4) {
                    let r = rightsMap['edit'];
                    let icon = r.icon || 'ti ti-edit';
                    buttons.push(`<a href="#" class="btn-grid-row-action text-primary" title="${r.title || 'Edit Designation'}" data-action="edit" data-url="${r.url || ''}" data-title="${r.title || 'Edit Designation'}" data-id="${targetId}" data-name="${titleValue}"><i class="${icon} fs-16"></i></a>`);
                }

                // Delete
                if (rightsMap['delete'] && (procStatus === 0 || procStatus === 2 || procStatus === 3) && actStatus !== 4) {
                    let r = rightsMap['delete'];
                    let icon = r.icon || 'ti ti-trash';
                    buttons.push(`<a href="#" class="btn-grid-row-action text-danger" title="${r.title || 'Delete Designation'}" data-action="delete" data-url="${r.url || ''}" data-title="${r.title || 'Delete Designation'}" data-id="${targetId}" data-name="${titleValue}"><i class="${icon} fs-16"></i></a>`);
                }

                // Sent To Checker
                if (rightsMap['stc'] && (procStatus === 0 || procStatus === 2) && actStatus !== 5 && actStatus !== 6) {
                    let r = rightsMap['stc'];
                    let icon = r.icon || 'ti ti-location';
                    buttons.push(`<a href="#" class="btn-grid-row-action text-warning" title="${r.title || 'Sent To Checker'}" data-action="stc" data-url="${r.url || ''}" data-title="${r.title || 'Sent To Checker'}" data-id="${targetId}" data-name="${titleValue}"><i class="${icon} fs-16"></i></a>`);
                }

                // Verify
                let canVerify = (typeof window.duty_carry_out[full.checker] !== 'undefined' || 1 === window.super_admin_status);
                if (rightsMap['verify'] && procStatus === 1 && canVerify) {
                    let r = rightsMap['verify'];
                    let icon = r.icon || 'ti ti-checks';
                    buttons.push(`<a href="#" class="btn-grid-row-action text-success" title="${r.title || 'Verify Designation'}" data-action="verify" data-url="${r.url || ''}" data-title="${r.title || 'Verify Designation'}" data-id="${targetId}" data-name="${titleValue}"><i class="${icon} fs-16"></i></a>`);
                }

                // Active
                if (rightsMap['active'] && (procStatus === 3 || procStatus === 2 || actStatus === 5 || procStatus === 0) && recStatus === 0) {
                    let r = rightsMap['active'];
                    let icon = r.icon || 'ti ti-creative-commons-sa';
                    buttons.push(`<a href="#" class="btn-grid-row-action text-success" title="${r.title || 'Activate Designation'}" data-action="active" data-url="${r.url || ''}" data-title="${r.title || 'Activate Designation'}" data-id="${targetId}" data-name="${titleValue}"><i class="${icon} fs-16"></i></a>`);
                }

                // Deactivate
                if (rightsMap['deactivate'] && (procStatus === 3 || procStatus === 2 || actStatus === 6 || procStatus === 0) && recStatus === 1) {
                    let r = rightsMap['deactivate'];
                    let icon = r.icon || 'ti ti-assembly-off';
                    buttons.push(`<a href="#" class="btn-grid-row-action text-dark" title="${r.title || 'Deactivate Designation'}" data-action="deactivate" data-url="${r.url || ''}" data-title="${r.title || 'Deactivate Designation'}" data-id="${targetId}" data-name="${titleValue}"><i class="${icon} fs-16"></i></a>`);
                }

                if (buttons.length === 0) {
                    return '<span class="text-muted fs-12">-</span>';
                }

                return `<div class="d-flex align-items-center gap-2 h-100 ps-1">${buttons.join('')}</div>`;
            }
        };
    }

    function initDesignationGrid() {
        let rightsList = getGridRights();
        let actionsColumnDef = createActionsColumnDef(rightsList);

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
                    let name = (params.data && (params.data.name || params.data.designation_name)) || '-';
                    return `<span class="fw-semibold text-dark"><i class="ti ti-briefcase me-1 text-primary"></i>${name}</span>`;
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
                    let val = Number(params.value);
                    let actLabel = window.actual_status[val] || 'New Add';
                    let badgeClass = 'bg-primary-transparent text-primary';
                    if (val === 4) badgeClass = 'bg-danger-transparent text-danger';
                    else if (val === 5) badgeClass = 'bg-dark-transparent text-dark';
                    else if (val === 6) badgeClass = 'bg-success-transparent text-success';
                    else if (val === 2) badgeClass = 'bg-warning-transparent text-warning';
                    else if (val === 3) badgeClass = 'bg-info-transparent text-info';
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
                    let st = Number(params.value);
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
        openMenu('/reference/designation/add', 'Add Designation', '', 'add', 1, 800);
    }

    function openEditModal(targetId) {
        currentDesignationUuid = targetId;
        openMenu('/reference/designation/edit', 'Edit Designation', targetId, 'edit', 1, 800);
    }

    function openViewModal(targetId) {
        currentDesignationUuid = targetId;
        openMenu('/reference/designation/view', 'Designation Details', targetId, 'view', 1, 800);
    }

    function openDeleteModal(targetId, name) {
        currentDesignationUuid = targetId;
        openMenu('/reference/designation/delete', 'Delete Designation', targetId, 'delete', 1, 600);
    }

    function openStcModal(targetId, name) {
        currentDesignationUuid = targetId;
        openMenu('/reference/designation/stc', 'Sent To Checker', targetId, 'stc', 1, 800);
    }

    function openVerifyModal(targetId, name) {
        currentDesignationUuid = targetId;
        openMenu('/reference/designation/verify', 'Verify Designation', targetId, 'verify', 1, 800);
    }

    function openStatusModal(actionType, targetId, name) {
        currentDesignationUuid = targetId;
        if (actionType === 'active') {
            openMenu('/reference/designation/active', 'Active Designation', targetId, 'active', 1, 600);
        } else {
            openMenu('/reference/designation/deactivate', 'Deactivate Designation', targetId, 'deactivate', 1, 600);
        }
    }

    function openHistoryModal(targetId, name) {
        targetId = targetId || currentDesignationUuid;
        if (!targetId) {
            showNotify('warning', 'Please select a designation first');
            return;
        }
        currentDesignationUuid = targetId;
        openMenu('/reference/designation/history', 'Designation Audit History', targetId, 'history', 1, 1140);
    }

    function openActivityModal(targetId) {
        targetId = targetId || currentDesignationUuid;
        if (!targetId) {
            showNotify('warning', 'Please select a designation first');
            return;
        }
        currentDesignationUuid = targetId;
        openMenu('/reference/designation/activity', 'Designation Activity Logs', targetId, 'activity', 1, 1140);
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
            openHistoryModal(uuid, 'Designation History');
        });

        // Activity & History buttons on top right of all sub-modals
        $(document).off('click.designationActivity', '.btn-open-activity').on('click.designationActivity', '.btn-open-activity', function (e) {
            e.preventDefault();
            let $actModal = $('#modal_designation_activity');
            if ($actModal.length && $actModal.find('tbody tr').length > 0) {
                if ($actModal.hasClass('show')) {
                    let $currentTop = $('.modal.show').last();
                    if ($currentTop[0] !== $actModal[0]) {
                        $currentTop.modal('hide');
                    }
                } else {
                    $actModal.modal('show');
                }
            } else {
                openActivityModal(currentDesignationUuid);
            }
        });

        $(document).off('click.designationHistory', '.btn-open-history').on('click.designationHistory', '.btn-open-history', function (e) {
            e.preventDefault();
            let $histModal = $('#modal_designation_history');
            if ($histModal.length && $histModal.find('tbody tr').length > 0) {
                if ($histModal.hasClass('show')) {
                    let $currentTop = $('.modal.show').last();
                    if ($currentTop[0] !== $histModal[0]) {
                        $currentTop.modal('hide');
                    }
                } else {
                    $histModal.modal('show');
                }
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
                if ($("#return_reason").val() && $("#return_reason").val().trim().length > 0) {
                    form_validation({control: "return_reason#text#mandatory", msg_text: "Return Reason"});
                }
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
        let $btn = $("#btn_save_designation_add");
        if ($btn.prop("disabled") || $btn.hasClass("submitting")) return;

        let status = form_validation({control: "name#text#mandatory*code#text#mandatory*title#text#mandatory*details#text#nonmandatory", msg_text: "Name*Code*Title*Details"});

        if (status === false) {
            showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
            return false;
        }

        $btn.addClass("submitting");
        let confirmed = await confirmDialog("Designation", "Are you sure you want to save the Designation?");

        if (!confirmed) {
            $btn.removeClass("submitting");
            return;
        }

        $btn.prop("disabled", true).addClass("disabled");

        let formData = new FormData();
        formData.append("name", $("#name").val());
        formData.append("code", $("#code").val());
        formData.append("title", $("#title").val());
        formData.append("details", $("#details").val());

        try {
            let response = await web.post('/reference/designation/add', formData);

            if (response.status === 200 || response.status === 201) {
                showAlertBox({title: "Success", message: "Saved Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
                closeModal('#modal_designation_add');
                refreshGrid();
            }
        } catch (err) {
            handleAjaxError(err, "Failed to save designation");
        } finally {
            $btn.prop("disabled", false).removeClass("disabled submitting");
        }
    }

    async function saveDesignationEdit() {
        let $btn = $("#btn_update_designation");
        if ($btn.prop("disabled") || $btn.hasClass("submitting")) return;

        let status = form_validation({control: "name#text#mandatory*code#text#mandatory*title#text#mandatory*details#text#nonmandatory", msg_text: "Name*Code*Title*Details"});

        if (status === false) {
            showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
            return false;
        }

        $btn.addClass("submitting");
        let confirmed = await confirmDialog("Designation", "Are you sure you want to update the Designation?");

        if (!confirmed) {
            $btn.removeClass("submitting");
            return;
        }

        $btn.prop("disabled", true).addClass("disabled");

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
                closeModal('#modal_designation_edit');
                refreshGrid();
            }
        } catch (err) {
            handleAjaxError(err, "Failed to update designation");
            refreshGrid();
        } finally {
            $btn.prop("disabled", false).removeClass("disabled submitting");
        }
    }

    async function saveDesignationDelete() {
        let $btn = $("#btn_confirm_delete");
        if ($btn.prop("disabled") || $btn.hasClass("submitting")) return;

        let status = form_validation({control: "delete_reason#text#mandatory", msg_text: "Delete Reason"});

        if (status === false) {
            showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
            return false;
        }

        $btn.addClass("submitting");
        let confirmed = await confirmDialog("Delete Designation", "Are you sure you want to delete this Designation?");

        if (!confirmed) {
            $btn.removeClass("submitting");
            return;
        }

        $btn.prop("disabled", true).addClass("disabled");

        let uuid = $('#delete_target_id').val() || currentDesignationUuid;

        let formData = new FormData();
        formData.append("delete_reason", $("#delete_reason").val());

        try {
            let response = await web.post(`/reference/designation/delete/${uuid}`, formData);

            if (response.status === 200) {
                showAlertBox({title: "Success", message: "Deleted Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
                closeModal('#modal_designation_delete');
                refreshGrid();
            }
        } catch (err) {
            handleAjaxError(err, "Failed to delete designation");
            refreshGrid();
        } finally {
            $btn.prop("disabled", false).removeClass("disabled submitting");
        }
    }

    async function saveDesignationStc() {
        let $btn = $("#btn_confirm_stc");
        if ($btn.prop("disabled") || $btn.hasClass("submitting")) return;

        let status = form_validation({control: "checker#select#mandatory*remarks#text#nonmandatory", msg_text: "Checker User*Remarks"});

        if (status === false) {
            showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
            return false;
        }

        $btn.addClass("submitting");
        let confirmed = await confirmDialog("Sent To Checker", "Are you sure you want to send this Designation to Checker?");

        if (!confirmed) {
            $btn.removeClass("submitting");
            return;
        }

        $btn.prop("disabled", true).addClass("disabled");

        let uuid = $('#stc_target_id').val() || currentDesignationUuid;

        let formData = new FormData();
        formData.append("checker", $("#checker").val());
        formData.append("remarks", $("#remarks").val());

        try {
            let response = await web.post(`/reference/designation/stc/${uuid}`, formData);

            if (response.status === 200) {
                showAlertBox({title: "Success", message: "Sent To Checker Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
                closeModal('#modal_designation_stc');
                refreshGrid();
            }
        } catch (err) {
            handleAjaxError(err, "Failed to send designation to checker");
            refreshGrid();
        } finally {
            $btn.prop("disabled", false).removeClass("disabled submitting");
        }
    }

    async function saveDesignationVerify(verifyType) {
        let $btn = verifyType === 0 ? $("#btn_verify_return") : $("#btn_verify_approve");
        if ($btn.prop("disabled") || $btn.hasClass("submitting")) return;

        if (verifyType === 0) {
            let status = form_validation({control: "return_reason#text#mandatory", msg_text: "Return Reason"});

            if (status === false) {
                showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
                return false;
            }

            $btn.addClass("submitting");
            let confirmed = await confirmDialog("Return Designation", "Are you sure you want to return this Designation?");

            if (!confirmed) {
                $btn.removeClass("submitting");
                return;
            }

            $btn.prop("disabled", true).addClass("disabled");

            let uuid = $('#verify_target_id').val() || currentDesignationUuid;

            let formData = new FormData();
            formData.append("verify_type", "0");
            formData.append("return_reason", $("#return_reason").val());

            try {
                let response = await web.post(`/reference/designation/verify/${uuid}`, formData);

                if (response.status === 200) {
                    showAlertBox({title: "Success", message: "Returned Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
                    closeModal('#modal_designation_verify');
                    refreshGrid();
                }
            } catch (err) {
                handleAjaxError(err, "Failed to return designation");
                refreshGrid();
            } finally {
                $btn.prop("disabled", false).removeClass("disabled submitting");
            }

        } else {
            $btn.addClass("submitting");
            let confirmed = await confirmDialog("Verify Designation", "Are you sure you want to approve this Designation?");

            if (!confirmed) {
                $btn.removeClass("submitting");
                return;
            }

            $btn.prop("disabled", true).addClass("disabled");

            let uuid = $('#verify_target_id').val() || currentDesignationUuid;

            let formData = new FormData();
            formData.append("verify_type", "1");

            try {
                let response = await web.post(`/reference/designation/verify/${uuid}`, formData);

                if (response.status === 200) {
                    showAlertBox({title: "Success", message: "Approved Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
                    closeModal('#modal_designation_verify');
                    refreshGrid();
                }
            } catch (err) {
                handleAjaxError(err, "Failed to approve designation");
                refreshGrid();
            } finally {
                $btn.prop("disabled", false).removeClass("disabled submitting");
            }
        }
    }

    async function saveDesignationActive() {
        let $btn = $("#btn_confirm_active");
        if ($btn.prop("disabled") || $btn.hasClass("submitting")) return;

        let status = form_validation({control: "active_reason#text#mandatory", msg_text: "Active Reason"});

        if (status === false) {
            showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
            return false;
        }

        $btn.addClass("submitting");
        let confirmed = await confirmDialog("Active Designation", "Are you sure you want to activate this Designation?");

        if (!confirmed) {
            $btn.removeClass("submitting");
            return;
        }

        $btn.prop("disabled", true).addClass("disabled");

        let uuid = $('#active_target_id').val() || currentDesignationUuid;

        let formData = new FormData();
        formData.append("active_reason", $("#active_reason").val());

        try {
            let response = await web.post(`/reference/designation/active/${uuid}`, formData);

            if (response.status === 200) {
                showAlertBox({title: "Success", message: "Activated Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
                closeModal('#modal_designation_active');
                refreshGrid();
            }
        } catch (err) {
            handleAjaxError(err, "Failed to activate designation");
            refreshGrid();
        } finally {
            $btn.prop("disabled", false).removeClass("disabled submitting");
        }
    }

    async function saveDesignationDeactivate() {
        let $btn = $("#btn_confirm_deactive");
        if ($btn.prop("disabled") || $btn.hasClass("submitting")) return;

        let status = form_validation({control: "deactivate_reason#text#mandatory", msg_text: "Deactivate Reason"});

        if (status === false) {
            showAlertBox({title: "Error", message: "Validation failed", type: "danger", sound: "/assets/sounds/error.mp3"});
            return false;
        }

        $btn.addClass("submitting");
        let confirmed = await confirmDialog("Deactivate Designation", "Are you sure you want to deactivate this Designation?");

        if (!confirmed) {
            $btn.removeClass("submitting");
            return;
        }

        $btn.prop("disabled", true).addClass("disabled");

        let uuid = $('#deactive_target_id').val() || currentDesignationUuid;

        let formData = new FormData();
        formData.append("deactivate_reason", $("#deactivate_reason").val());

        try {
            let response = await web.post(`/reference/designation/deactivate/${uuid}`, formData);

            if (response.status === 200) {
                showAlertBox({title: "Success", message: "Deactivated Successfully", type: "info", sound: "/assets/sounds/success.mp3"});
                closeModal('#modal_designation_deactive');
                refreshGrid();
            }
        } catch (err) {
            handleAjaxError(err, "Failed to deactivate designation");
            refreshGrid();
        } finally {
            $btn.prop("disabled", false).removeClass("disabled submitting");
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

        $(document).on("submit.designationForm", "#form_designation_verify", function (e) {
            e.preventDefault();
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
