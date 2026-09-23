'use strict';

(function ($) {
    let designationGrid = null;

    $(document).ready(function () {
        initDesignationGrid();
        bindEvents();
    });

    function initDesignationGrid() {
        designationGrid = initAgDataGrid({
            gridContainer: '#designation_ag_grid',
            url: '/reference/designation/list',
            exportFileName: 'designations_list',
            defaultOrderCol: 'id',
            defaultOrderDir: 'desc',
            totalBadgeSelector: '#designation_total_badge',
            columnDefs: [
                {
                    headerName: '#',
                    isIndexCol: true,
                    width: 70,
                    minWidth: 65
                },
                {
                    field: 'name',
                    colId: 'name',
                    headerName: 'Designation Name',
                    width: 250,
                    minWidth: 180,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        return `<span class="fw-semibold text-dark"><i class="ti ti-briefcase me-1 text-primary"></i>${params.value || '-'}</span>`;
                    }
                },
                {
                    field: 'title',
                    colId: 'title',
                    headerName: 'Title',
                    width: 200,
                    minWidth: 150,
                    sortable: true,
                    filter: true
                },
                {
                    field: 'code',
                    colId: 'code',
                    headerName: 'Code',
                    width: 140,
                    minWidth: 110,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        return params.value ? `<span class="badge bg-light text-dark border fs-12">${params.value}</span>` : '<span class="text-muted">-</span>';
                    }
                },
                {
                    field: 'status_text',
                    colId: 'status',
                    headerName: 'Status',
                    width: 130,
                    minWidth: 110,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        let row = params.data || {};
                        if (row.status === 1) {
                            return `<span class="badge bg-success-transparent text-success fw-semibold"><i class="ti ti-check me-1"></i>Active</span>`;
                        }
                        return `<span class="badge bg-danger-transparent text-danger fw-semibold"><i class="ti ti-x me-1"></i>Inactive</span>`;
                    }
                },
                {
                    field: 'checker_text',
                    colId: 'checker',
                    headerName: 'Checker Status',
                    width: 160,
                    minWidth: 130,
                    sortable: true,
                    cellRenderer: function (params) {
                        let row = params.data || {};
                        if (row.checker === 2) {
                            return `<span class="badge bg-primary-transparent text-primary fw-semibold"><i class="ti ti-circle-check me-1"></i>Verified</span>`;
                        } else if (row.checker === 1) {
                            return `<span class="badge bg-warning-transparent text-warning fw-semibold"><i class="ti ti-clock me-1"></i>Sent to Checker</span>`;
                        }
                        return `<span class="badge bg-secondary-transparent text-secondary fw-semibold">Pending</span>`;
                    }
                },
                {
                    field: 'created_at',
                    colId: 'created_at',
                    headerName: 'Created At',
                    width: 180,
                    minWidth: 160,
                    sortable: true,
                    cellRenderer: function (params) {
                        return `<span class="text-muted fs-12"><i class="ti ti-calendar me-1"></i>${params.value || '-'}</span>`;
                    }
                },
                {
                    headerName: 'Actions',
                    width: 180,
                    minWidth: 160,
                    pinned: 'right',
                    sortable: false,
                    filter: false,
                    cellRenderer: function (params) {
                        let row = params.data;
                        if (!row) return '';
                        let statusBtn = row.status === 1
                            ? `<button type="button" class="btn btn-sm btn-icon btn-light text-warning btn-action-deactivate" data-id="${row.id}" title="Deactivate"><i class="ti ti-ban fs-14"></i></button>`
                            : `<button type="button" class="btn btn-sm btn-icon btn-light text-success btn-action-active" data-id="${row.id}" title="Activate"><i class="ti ti-check fs-14"></i></button>`;

                        return `
                            <div class="d-flex align-items-center gap-1">
                                <button type="button" class="btn btn-sm btn-icon btn-light text-info btn-action-edit" data-id="${row.id}" data-name="${encodeURIComponent(row.name || '')}" data-title="${encodeURIComponent(row.title || '')}" data-code="${encodeURIComponent(row.code || '')}" data-details="${encodeURIComponent(row.details || '')}" title="Edit">
                                    <i class="ti ti-edit fs-14"></i>
                                </button>
                                ${statusBtn}
                                <button type="button" class="btn btn-sm btn-icon btn-light text-primary btn-action-verify" data-id="${row.id}" title="Verify">
                                    <i class="ti ti-checks fs-14"></i>
                                </button>
                                <button type="button" class="btn btn-sm btn-icon btn-light text-danger btn-action-delete" data-id="${row.id}" title="Delete">
                                    <i class="ti ti-trash fs-14"></i>
                                </button>
                            </div>
                        `;
                    }
                }
            ]
        });
    }

    function bindEvents() {
        $('#btn_open_add_modal').on('click', function () {
            $('#form_designation')[0].reset();
            $('#designation_id').val('');
            $('#modal_designation_title').text('Add Designation');
            $('#modal_designation_form').modal('show');
        });

        $('#form_designation').on('submit', function (e) {
            e.preventDefault();
            let id = $('#designation_id').val();
            let url = id ? `/reference/designation/edit/${id}` : '/reference/designation/add';
            let formData = $(this).serialize();

            $.post(url, formData)
                .done(function (res) {
                    $('#modal_designation_form').modal('hide');
                    if (window.toastNotification) {
                        toastNotification('success', res.message || 'Saved successfully');
                    }
                    if (designationGrid && designationGrid.gridApi) {
                        designationGrid.gridApi.refreshServerSide({ purge: true });
                    }
                })
                .fail(function (xhr) {
                    let msg = 'Failed to save designation';
                    if (xhr.responseJSON && xhr.responseJSON.message) {
                        msg = xhr.responseJSON.message;
                    }
                    if (window.toastNotification) {
                        toastNotification('danger', msg);
                    } else {
                        alert(msg);
                    }
                });
        });

        $(document).on('click', '.btn-action-edit', function () {
            let id = $(this).data('id');
            $('#designation_id').val(id);
            $('#designation_name').val(decodeURIComponent($(this).data('name') || ''));
            $('#designation_title').val(decodeURIComponent($(this).data('title') || ''));
            $('#designation_code').val(decodeURIComponent($(this).data('code') || ''));
            $('#designation_details').val(decodeURIComponent($(this).data('details') || ''));
            $('#modal_designation_title').text('Edit Designation');
            $('#modal_designation_form').modal('show');
        });

        $(document).on('click', '.btn-action-delete', function () {
            let id = $(this).data('id');
            if (!confirm('Are you sure you want to delete this designation?')) return;
            let csrfToken = $('input[name="_csrf"]').val();

            $.post(`/reference/designation/delete/${id}`, { _csrf: csrfToken })
                .done(function (res) {
                    if (window.toastNotification) {
                        toastNotification('success', res.message || 'Deleted successfully');
                    }
                    if (designationGrid && designationGrid.gridApi) {
                        designationGrid.gridApi.refreshServerSide({ purge: true });
                    }
                })
                .fail(function (xhr) {
                    let msg = (xhr.responseJSON && xhr.responseJSON.message) || 'Error deleting designation';
                    if (window.toastNotification) toastNotification('danger', msg);
                });
        });

        $(document).on('click', '.btn-action-active', function () {
            let id = $(this).data('id');
            let csrfToken = $('input[name="_csrf"]').val();
            $.post(`/reference/designation/active/${id}`, { _csrf: csrfToken }).done(function () {
                if (designationGrid && designationGrid.gridApi) {
                    designationGrid.gridApi.refreshServerSide({ purge: true });
                }
            });
        });

        $(document).on('click', '.btn-action-deactivate', function () {
            let id = $(this).data('id');
            let csrfToken = $('input[name="_csrf"]').val();
            $.post(`/reference/designation/deactivate/${id}`, { _csrf: csrfToken }).done(function () {
                if (designationGrid && designationGrid.gridApi) {
                    designationGrid.gridApi.refreshServerSide({ purge: true });
                }
            });
        });

        $(document).on('click', '.btn-action-verify', function () {
            let id = $(this).data('id');
            let csrfToken = $('input[name="_csrf"]').val();
            $.post(`/reference/designation/verify/${id}`, { _csrf: csrfToken }).done(function () {
                if (designationGrid && designationGrid.gridApi) {
                    designationGrid.gridApi.refreshServerSide({ purge: true });
                }
            });
        });
    }
})(jQuery);
