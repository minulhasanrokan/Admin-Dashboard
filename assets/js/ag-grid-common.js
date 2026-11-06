/**
 * ag-grid-common.js
 * Universal Reusable Enterprise AG Grid Framework for Admin Dashboard
 * Provides Server-Side Pagination, Column Filtering Popups, Sorting, Global Search, and Button Pagination.
 */
'use strict';

(function (window, $) {
    let activeColumnPopup = null;

    // Parse application datetime constants from CSP-compliant JSON script block
    function getAppDateTimeConstants() {
        if (window.APP_DATETIME_CONSTANTS) return window.APP_DATETIME_CONSTANTS;
        const scriptEl = document.getElementById('app_datetime_constants');
        if (scriptEl && scriptEl.textContent) {
            try {
                window.APP_DATETIME_CONSTANTS = JSON.parse(scriptEl.textContent.trim());
                return window.APP_DATETIME_CONSTANTS;
            } catch (e) {}
        }
        return {
            JsInputDateFormat: 'DD/MM/YYYY',
            JsInputDateTimeFormat: 'DD/MM/YYYY HH:mm',
            JsViewDateFormat: 'DD MMM YYYY',
            JsViewDateTimeFormat: 'DD MMM YYYY hh:mm A'
        };
    }

    // Global listener to close popups on outside click
    $(document).on('click.agCommonPopup', function (e) {
        if (!$(e.target).closest('.jqx-col-popup, .ag-col-menu-trigger, .daterangepicker').length) {
            closeCommonColumnPopup();
        }
    });

    function closeCommonColumnPopup() {
        if (activeColumnPopup) {
            activeColumnPopup.remove();
            activeColumnPopup = null;
            $('.ag-col-menu-trigger').removeClass('active');
        }
    }

    function normalizeOperator(op) {
        if (!op) return 'contains';
        switch (op) {
            case 'lessThan': return 'less_than';
            case 'greaterThan': return 'greater_than';
            case 'lessThanOrEqual': return 'less_than_or_equal';
            case 'greaterThanOrEqual': return 'greater_than_or_equal';
            case 'startsWith': return 'starts_with';
            case 'endsWith': return 'ends_with';
            case 'notContains': return 'not_contains';
            case 'notEqual': return 'not_equal';
            case 'equals': return 'equals';
            default: return op;
        }
    }

    function escapeHtml(text) {
        if (!text) return "";
        return String(text)
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }

    /**
     * Initializes a common, reusable enterprise AG Grid instance.
     * @param {Object} config
     * @returns {Object} Grid Controller instance
     */
    function initAgDataGrid(config) {
        let gridContainer = typeof config.gridContainer === 'string'
            ? document.querySelector(config.gridContainer)
            : config.gridContainer;

        if (!gridContainer) {
            console.error('initAgDataGrid: Target container not found:', config.gridContainer);
            return null;
        }

        let url = config.url || '';
        let exportFileName = config.exportFileName || 'grid_export';
        let defaultOrderCol = config.defaultOrderCol || 'id';
        let defaultOrderDir = config.defaultOrderDir || 'desc';
        let initialPageSize = config.pageSize || 10;
        let searchSelector = config.searchSelector || '#ag_quick_filter';
        let searchClearSelector = config.searchClearSelector || '#btn_ag_search_clear';
        let refreshSelector = config.refreshSelector || '#btn_refresh_ag_grid';
        let exportSelector = config.exportSelector || '#btn_ag_export_csv';
        let paginationSelector = config.paginationSelector || '#ag_custom_pagination';
        let totalBadgeSelector = config.totalBadgeSelector || '#backup_total_badge';

        let currentPage = 1;
        let currentPageSize = initialPageSize;
        let totalRecords = 0;
        let currentOrderCol = defaultOrderCol;
        let currentOrderDir = defaultOrderDir;
        let columnFilterStates = {}; // Stores { op1, val1, logic, op2, val2 } per colId
        let searchDebounceTimer = null;
        let gridApi = null;
        let headerInstances = [];

        function updateAllSortIndicators() {
            headerInstances.forEach(h => {
                if (h && typeof h.updateSortIndicator === 'function') {
                    h.updateSortIndicator();
                }
            });
        }

        // Custom Header Component with jqx-style popover
        class CommonGridHeader {
            init(params) {
                this.params = params;
                this.eGui = document.createElement('div');

                headerInstances.push(this);

                let colDef = params.column.getColDef();
                let title = params.displayName;
                let isFilterable = colDef.filter !== false;
                let isSortable = colDef.sortable !== false;
                let colId = params.column.getColId();
                let hasFilter = !!columnFilterStates[colId];

                let titleContent = escapeHtml(title);
                let titleAttr = escapeHtml(title);

                if (colDef.headerIconHtml) {
                    titleContent = colDef.headerIconHtml;
                    if (colDef.headerIconTitle) {
                        titleAttr = escapeHtml(colDef.headerIconTitle);
                    }
                } else if (colDef.headerIcon) {
                    let iconStr = colDef.headerIcon;
                    if (iconStr.indexOf('/') !== -1 || iconStr.indexOf('.') !== -1) {
                        let src = iconStr.startsWith('http') || iconStr.startsWith('/') ? iconStr : `/assets/images/${iconStr}`;
                        titleContent = `<img src="${src}" alt="${titleAttr}" style="width:16px;height:16px;">`;
                    } else {
                        titleContent = `<i class="${iconStr} fs-15"></i>`;
                    }
                    if (colDef.headerIconTitle) {
                        titleAttr = escapeHtml(colDef.headerIconTitle);
                    }
                }

                let isCentered = (colDef.headerClass && String(colDef.headerClass).indexOf('text-center') !== -1) ||
                                 (colDef.cellClass && String(colDef.cellClass).indexOf('justify-content-center') !== -1);
                let justifyClass = isCentered ? 'justify-content-center' : 'justify-content-between';

                this.eGui.className = `ag-custom-header-wrapper d-flex align-items-center ${justifyClass} w-100`;

                let headerHtml = `
                    <div class="ag-header-cell-label d-flex align-items-center gap-1 overflow-hidden" role="presentation" title="${titleAttr}">
                        <span class="ag-header-cell-text text-truncate">${titleContent}</span>
                        ${isSortable ? '<span class="ag-header-icon ag-sort-indicator ms-1 ag-sort-indicator-w"></span>' : ''}
                    </div>
                `;

                if (isFilterable || isSortable) {
                    headerHtml += `
                        <div class="ag-col-menu-trigger ${hasFilter ? 'has-filter' : ''}" title="Sort & Filter Column">
                            <i class="ti ti-menu-2 fs-12"></i>
                        </div>
                    `;
                }

                this.eGui.innerHTML = headerHtml;
                this.sortIndicator = this.eGui.querySelector('.ag-sort-indicator');
                this.menuBtn = this.eGui.querySelector('.ag-col-menu-trigger');
                this.labelEl = this.eGui.querySelector('.ag-header-cell-label');

                if (isSortable && this.labelEl) {
                    this.labelEl.style.cursor = 'pointer';
                    this.labelEl.addEventListener('click', (e) => {
                        e.stopPropagation();
                        if (currentOrderCol === colId) {
                            if (currentOrderDir === 'asc') {
                                currentOrderDir = 'desc';
                            } else if (currentOrderDir === 'desc') {
                                currentOrderCol = defaultOrderCol;
                                currentOrderDir = defaultOrderDir;
                            } else {
                                currentOrderDir = 'asc';
                            }
                        } else {
                            currentOrderCol = colId;
                            currentOrderDir = 'asc';
                        }

                        updateAllSortIndicators();
                        controller.reload(1, currentPageSize);
                    });
                }

                if (this.menuBtn) {
                    this.menuBtn.addEventListener('click', (e) => {
                        e.stopPropagation();
                        e.preventDefault();
                        showColumnMenuPopup(this.menuBtn, params);
                    });
                }

                this.updateSortIndicator();
            }

            updateSortIndicator() {
                if (!this.sortIndicator) return;
                let colId = this.params.column.getColId();
                let isDefaultSort = (currentOrderCol === defaultOrderCol && currentOrderDir === defaultOrderDir);

                if (!isDefaultSort && currentOrderCol === colId) {
                    if (currentOrderDir === 'asc') {
                        this.sortIndicator.innerHTML = '<i class="ti ti-arrow-up text-primary fs-11"></i>';
                    } else if (currentOrderDir === 'desc') {
                        this.sortIndicator.innerHTML = '<i class="ti ti-arrow-down text-primary fs-11"></i>';
                    } else {
                        this.sortIndicator.innerHTML = '';
                    }
                } else {
                    this.sortIndicator.innerHTML = '';
                }
            }

            getGui() {
                return this.eGui;
            }

            destroy() {
                headerInstances = headerInstances.filter(h => h !== this);
            }
        }

        function showColumnMenuPopup(buttonElement, params) {
            closeCommonColumnPopup();

            let colDef = params.column.getColDef();
            let colId = params.column.getColId();
            let filterType = colDef.filterType || 'string'; // 'string', 'date', 'size', 'number', 'select'
            let currentFilter = columnFilterStates[colId] || { op1: '', val1: '', logic: 'AND', op2: '', val2: '' };

            $(buttonElement).addClass('active');

            let opOptions = '';
            let labelType = 'text';
            let isDate = false;
            let isSelect = false;

            if (filterType === 'date') {
                labelType = 'date';
                isDate = true;
                opOptions = `
                    <option value="equals" ${currentFilter.op1 === 'equals' ? 'selected' : ''}>equal</option>
                    <option value="notEqual" ${currentFilter.op1 === 'notEqual' ? 'selected' : ''}>not equal</option>
                    <option value="lessThan" ${currentFilter.op1 === 'lessThan' || !currentFilter.op1 ? 'selected' : ''}>less than</option>
                    <option value="greaterThan" ${currentFilter.op1 === 'greaterThan' ? 'selected' : ''}>greater than</option>
                `;
            } else if (filterType === 'select' || (Array.isArray(colDef.filterOptions) && colDef.filterOptions.length > 0)) {
                labelType = 'option';
                isSelect = true;
                opOptions = `
                    <option value="equals" ${currentFilter.op1 === 'equals' || !currentFilter.op1 ? 'selected' : ''}>equals</option>
                    <option value="notEqual" ${currentFilter.op1 === 'notEqual' ? 'selected' : ''}>not equal</option>
                `;
            } else if (filterType === 'size') {
                labelType = 'file size';
                opOptions = `
                    <option value="lessThan" ${currentFilter.op1 === 'lessThan' || !currentFilter.op1 ? 'selected' : ''}>less than</option>
                    <option value="greaterThan" ${currentFilter.op1 === 'greaterThan' ? 'selected' : ''}>greater than</option>
                    <option value="lessThanOrEqual" ${currentFilter.op1 === 'lessThanOrEqual' ? 'selected' : ''}>less than or equal</option>
                    <option value="greaterThanOrEqual" ${currentFilter.op1 === 'greaterThanOrEqual' ? 'selected' : ''}>greater than or equal</option>
                    <option value="equals" ${currentFilter.op1 === 'equals' ? 'selected' : ''}>equals</option>
                    <option value="notEqual" ${currentFilter.op1 === 'notEqual' ? 'selected' : ''}>not equal</option>
                    <option value="contains" ${currentFilter.op1 === 'contains' ? 'selected' : ''}>contains</option>
                `;
            } else if (filterType === 'number') {
                labelType = 'number';
                opOptions = `
                    <option value="equals" ${currentFilter.op1 === 'equals' ? 'selected' : ''}>equal</option>
                    <option value="notEqual" ${currentFilter.op1 === 'notEqual' ? 'selected' : ''}>not equal</option>
                    <option value="lessThan" ${currentFilter.op1 === 'lessThan' || !currentFilter.op1 ? 'selected' : ''}>less than</option>
                    <option value="greaterThan" ${currentFilter.op1 === 'greaterThan' ? 'selected' : ''}>greater than</option>
                    <option value="lessThanOrEqual" ${currentFilter.op1 === 'lessThanOrEqual' ? 'selected' : ''}>less than or equal</option>
                    <option value="greaterThanOrEqual" ${currentFilter.op1 === 'greaterThanOrEqual' ? 'selected' : ''}>greater than or equal</option>
                `;
            } else {
                labelType = 'text';
                opOptions = `
                    <option value="contains" ${currentFilter.op1 === 'contains' || !currentFilter.op1 ? 'selected' : ''}>contains</option>
                    <option value="notContains" ${currentFilter.op1 === 'notContains' ? 'selected' : ''}>does not contain</option>
                    <option value="equals" ${currentFilter.op1 === 'equals' ? 'selected' : ''}>equals</option>
                    <option value="notEqual" ${currentFilter.op1 === 'notEqual' ? 'selected' : ''}>not equal</option>
                    <option value="startsWith" ${currentFilter.op1 === 'startsWith' ? 'selected' : ''}>starts with</option>
                    <option value="endsWith" ${currentFilter.op1 === 'endsWith' ? 'selected' : ''}>ends with</option>
                `;
            }

            let op2Options = opOptions.replace(`value="${currentFilter.op2}"`, `value="${currentFilter.op2}" selected`);
            let dynamicJsDateFormat = getAppDateTimeConstants().JsInputDateFormat || 'DD/MM/YYYY';

            function renderValueInput(val, inputClass) {
                if (isSelect) {
                    let optionsHtml = `<option value="">-- Select Option --</option>`;
                    (colDef.filterOptions || []).forEach(opt => {
                        let optVal = (typeof opt === 'object' && opt !== null) ? opt.value : opt;
                        let optLabel = (typeof opt === 'object' && opt !== null) ? opt.label : opt;
                        let isSel = String(val) === String(optVal) ? 'selected' : '';
                        optionsHtml += `<option value="${escapeHtml(optVal)}" ${isSel}>${escapeHtml(optLabel)}</option>`;
                    });
                    return `<select class="form-select form-select-sm ${inputClass} fs-12">${optionsHtml}</select>`;
                }
                let placeholderText = 'Value...';
                if (isDate) {
                    placeholderText = dynamicJsDateFormat;
                } else if (filterType === 'size') {
                    placeholderText = 'e.g. 10 MB, 500 KB...';
                }
                return `<input type="text" class="form-control form-control-sm ${inputClass} fs-12 ${isDate ? 'jqx-date-input' : ''}" placeholder="${placeholderText}" value="${escapeHtml(val || '')}">`;
            }

            let isFilterable = colDef.filter !== false;
            let isSortable = colDef.sortable !== false;

            let sortSectionHtml = '';
            if (isSortable) {
                sortSectionHtml = `
                    <!-- Sort Actions -->
                    <div class="jqx-sort-option" data-action="sort-asc">
                        <span class="text-primary fw-bold jqx-sort-icon">A&darr;</span>
                        <span>Sort Ascending</span>
                    </div>
                    <div class="jqx-sort-option" data-action="sort-desc">
                        <span class="text-primary fw-bold jqx-sort-icon">Z&uarr;</span>
                        <span>Sort Descending</span>
                    </div>
                    <div class="jqx-sort-option" data-action="sort-clear">
                        <span class="text-secondary fw-bold jqx-sort-icon">E&times;</span>
                        <span>Remove Sort</span>
                    </div>
                `;
            }

            let filterSectionHtml = '';
            if (isFilterable) {
                filterSectionHtml = `
                    ${isSortable ? '<hr class="my-2 border-secondary-subtle">' : ''}
                    <!-- Filter Conditions -->
                    <div class="fs-12 text-secondary mb-1">Show rows where ${labelType}:</div>

                    <!-- Condition 1 -->
                    <div class="mb-2">
                        <select class="form-select form-select-sm jqx-op1 mb-1 fs-12">
                            ${opOptions}
                        </select>
                        <div class="input-group input-group-sm">
                            ${renderValueInput(currentFilter.val1, 'jqx-val1')}
                        </div>
                    </div>

                    <!-- Logic AND / OR -->
                    <div class="mb-2">
                        <select class="form-select form-select-sm jqx-logic fs-12 jqx-logic-w">
                            <option value="AND" ${currentFilter.logic === 'AND' ? 'selected' : ''}>And</option>
                            <option value="OR" ${currentFilter.logic === 'OR' ? 'selected' : ''}>Or</option>
                        </select>
                    </div>

                    <!-- Condition 2 -->
                    <div class="mb-2">
                        <select class="form-select form-select-sm jqx-op2 mb-1 fs-12">
                            ${op2Options}
                        </select>
                        <div class="input-group input-group-sm">
                            ${renderValueInput(currentFilter.val2, 'jqx-val2')}
                        </div>
                    </div>

                    <!-- Button Actions -->
                    <div class="d-flex align-items-center justify-content-end gap-2 pt-2 border-top">
                        <button type="button" class="btn btn-sm btn-primary px-3 jqx-btn-filter fs-12">Filter</button>
                        <button type="button" class="btn btn-sm btn-light border px-3 jqx-btn-clear fs-12">Clear</button>
                    </div>
                `;
            }

            let popupHtml = `
                <div class="jqx-col-popup" id="jqx_active_popup">
                    ${sortSectionHtml}
                    ${filterSectionHtml}
                </div>
            `;

            $('body').append(popupHtml);
            activeColumnPopup = $('#jqx_active_popup');

            if (isDate && typeof $.fn.daterangepicker === 'function') {
                activeColumnPopup.find('.jqx-date-input').daterangepicker({
                    singleDatePicker: true,
                    showDropdowns: true,
                    autoUpdateInput: false,
                    locale: {
                        format: dynamicJsDateFormat
                    }
                }).on('apply.daterangepicker', function (ev, picker) {
                    $(this).val(picker.startDate.format(dynamicJsDateFormat));
                });
            }

            let rect = buttonElement.getBoundingClientRect();
            let top = rect.bottom + window.scrollY + 4;
            let left = rect.left + window.scrollX - 220;
            if (left < 10) left = 10;

            activeColumnPopup.css({ top: `${top}px`, left: `${left}px` });

            // Sort Action
            activeColumnPopup.on('click', '.jqx-sort-option', function () {
                let action = $(this).attr('data-action');
                if (action === 'sort-asc') {
                    currentOrderCol = colId;
                    currentOrderDir = 'asc';
                } else if (action === 'sort-desc') {
                    currentOrderCol = colId;
                    currentOrderDir = 'desc';
                } else if (action === 'sort-clear') {
                    currentOrderCol = defaultOrderCol;
                    currentOrderDir = defaultOrderDir;
                }
                closeCommonColumnPopup();
                updateAllSortIndicators();
                controller.reload(1, currentPageSize);
            });

            // Filter Action
            activeColumnPopup.on('click', '.jqx-btn-filter', function () {
                let op1 = activeColumnPopup.find('.jqx-op1').val();
                let val1 = activeColumnPopup.find('.jqx-val1').val().trim();
                let logic = activeColumnPopup.find('.jqx-logic').val();
                let op2 = activeColumnPopup.find('.jqx-op2').val();
                let val2 = activeColumnPopup.find('.jqx-val2').val().trim();

                if (!val1 && !val2) {
                    columnFilterStates[colId] = null;
                    $(buttonElement).removeClass('has-filter');
                } else {
                    columnFilterStates[colId] = { op1, val1, logic, op2, val2 };
                    $(buttonElement).addClass('has-filter');
                }

                closeCommonColumnPopup();
                controller.reload(1, currentPageSize);
            });

            // Clear Action
            activeColumnPopup.on('click', '.jqx-btn-clear', function () {
                columnFilterStates[colId] = null;
                $(buttonElement).removeClass('has-filter');
                closeCommonColumnPopup();
                controller.reload(1, currentPageSize);
            });
        }

        // Prepare column definitions
        let colDefs = (config.columnDefs || []).map(col => {
            if (col.isIndexCol) {
                return {
                    ...col,
                    headerName: col.headerName || '#',
                    colId: col.colId || 'row_index',
                    valueGetter: function (params) {
                        return (currentPage - 1) * currentPageSize + params.node.rowIndex + 1;
                    },
                    width: col.width || 70,
                    minWidth: col.minWidth || 45,
                    sortable: false,
                    filter: false,
                    cellClass: col.cellClass || 'd-flex align-items-center justify-content-center text-muted fw-medium'
                };
            }
            return {
                sortable: col.sortable !== false,
                filter: col.filter !== false,
                filterType: col.filterType || 'string',
                ...col
            };
        });

        const gridOptions = {
            components: {
                agColumnHeader: CommonGridHeader
            },
            rowHeight: 36,
            headerHeight: 36,
            suppressPaginationPanel: true,
            suppressHorizontalScroll: false,
            alwaysShowHorizontalScroll: false,
            alwaysShowVerticalScroll: false,
            animateRows: true,
            rowSelection: 'single',
            columnDefs: colDefs,
            defaultColDef: {
                sortable: true,
                resizable: true,
                filter: true,
                suppressMovable: false,
                unSortIcon: true,
                ...(config.defaultColDef || {})
            },
            overlayLoadingTemplate: `
                <div class="d-flex flex-column align-items-center justify-content-center p-4">
                    <div class="spinner-border text-primary mb-2 ag-overlay-spinner" role="status"></div>
                    <span class="fs-13 fw-medium text-muted">Loading records...</span>
                </div>
            `,
            overlayNoRowsTemplate: `
                <div class="d-flex flex-column align-items-center justify-content-center p-5">
                    <div class="avatar avatar-lg bg-light rounded-circle text-muted mb-2 d-flex align-items-center justify-content-center ag-overlay-avatar">
                        <i class="ti ti-database-off fs-24"></i>
                    </div>
                    <h6 class="fw-bold text-dark mb-1">No Records Found</h6>
                    <span class="text-muted fs-12">No data matches your current search or filter criteria.</span>
                </div>
            `,
            rowData: []
        };

        // Create AG Grid instance
        if (typeof agGrid !== 'undefined' && agGrid.createGrid) {
            gridApi = agGrid.createGrid(gridContainer, gridOptions);
        }

        // Render Pagination UI
        function renderPagination(page, total, limit) {
            let totalPages = Math.ceil(total / limit) || 1;
            if (page > totalPages) page = totalPages;

            let startRow = total === 0 ? 0 : (page - 1) * limit + 1;
            let endRow = Math.min(page * limit, total);

            $("#ag_page_start").text(startRow);
            $("#ag_page_end").text(endRow);
            $("#ag_page_total").text(total);

            let html = `
                <li class="page-item ${page <= 1 ? 'disabled' : ''}">
                    <a class="page-link" href="#" data-page="${page - 1}" aria-label="Previous" title="Previous Page">
                        <i class="ti ti-chevron-left fs-13"></i>
                    </a>
                </li>
            `;

            let maxVisible = 5;
            let startPage = 1;
            let endPage = totalPages;

            if (totalPages > maxVisible) {
                let half = Math.floor(maxVisible / 2);
                startPage = Math.max(1, page - half);
                endPage = Math.min(totalPages, startPage + maxVisible - 1);

                if (endPage - startPage < maxVisible - 1) {
                    startPage = Math.max(1, endPage - maxVisible + 1);
                }
            }

            if (startPage > 1) {
                html += `<li class="page-item"><a class="page-link" href="#" data-page="1">1</a></li>`;
                if (startPage > 2) {
                    html += `<li class="page-item disabled"><span class="page-link border-0 text-muted px-1">...</span></li>`;
                }
            }

            for (let i = startPage; i <= endPage; i++) {
                html += `
                    <li class="page-item ${i === page ? 'active' : ''}">
                        <a class="page-link" href="#" data-page="${i}">${i}</a>
                    </li>
                `;
            }

            if (endPage < totalPages) {
                if (endPage < totalPages - 1) {
                    html += `<li class="page-item disabled"><span class="page-link border-0 text-muted px-1">...</span></li>`;
                }
                html += `<li class="page-item"><a class="page-link" href="#" data-page="${totalPages}">${totalPages}</a></li>`;
            }

            html += `
                <li class="page-item ${page >= totalPages ? 'disabled' : ''}">
                    <a class="page-link" href="#" data-page="${page + 1}" aria-label="Next" title="Next Page">
                        <i class="ti ti-chevron-right fs-13"></i>
                    </a>
                </li>
            `;

            $("#ag_pagination_list").html(html);
        }

        // Controller API Object
        const controller = {
            gridApi: gridApi,

            reload: function (page, limit, callback) {
                if (!gridApi) return;

                currentPage = page || currentPage || 1;
                currentPageSize = limit || currentPageSize || 10;

                gridApi.showLoadingOverlay();

                let search = $(searchSelector).val() || "";

                let activeFilters = [];
                for (let field in columnFilterStates) {
                    let filterState = columnFilterStates[field];
                    if (!filterState) continue;

                    if (filterState.val1) {
                        activeFilters.push({
                            field: field,
                            operator: normalizeOperator(filterState.op1),
                            value: filterState.val1
                        });
                    }
                    if (filterState.val2) {
                        activeFilters.push({
                            field: field,
                            operator: normalizeOperator(filterState.op2),
                            value: filterState.val2
                        });
                    }
                }

                let formData = new FormData();
                formData.append("page", currentPage);
                formData.append("limit", currentPageSize);
                formData.append("search", search);
                formData.append("order_column", currentOrderCol);
                formData.append("order_dir", currentOrderDir);
                if (activeFilters.length > 0) {
                    formData.append("filters", JSON.stringify(activeFilters));
                }

                web.post(url, formData)
                    .then(function (response) {
                        let resPayload = response.data;
                        let payloadData = (resPayload && resPayload.data) ? resPayload.data : resPayload;
                        let items = [];

                        if (payloadData && Array.isArray(payloadData.items)) {
                            items = payloadData.items;
                        } else if (payloadData && Array.isArray(payloadData.data)) {
                            items = payloadData.data;
                        } else if (payloadData && Array.isArray(payloadData.rows)) {
                            items = payloadData.rows;
                        } else if (Array.isArray(payloadData)) {
                            items = payloadData;
                        }

                        totalRecords = (payloadData && typeof payloadData.total === 'number') ? payloadData.total : items.length;

                        $(totalBadgeSelector).text(`${totalRecords} ${totalRecords === 1 ? 'Record' : 'Records'}`);

                        gridApi.setGridOption('rowData', items);

                        if (items.length === 0) {
                            gridApi.showNoRowsOverlay();
                        } else {
                            gridApi.hideOverlay();
                        }

                        renderPagination(currentPage, totalRecords, currentPageSize);

                        if (typeof config.onRowDataLoaded === 'function') {
                            config.onRowDataLoaded(items, totalRecords);
                        }
                    })
                    .catch(function (error) {
                        console.error("CommonGrid Error loading server data:", error);
                        $(totalBadgeSelector).text("0 Records");
                        totalRecords = 0;
                        gridApi.setGridOption('rowData', []);
                        gridApi.showNoRowsOverlay();
                        renderPagination(1, 0, currentPageSize);
                    })
                    .finally(function () {
                        if (typeof callback === "function") callback();
                    });
            },

            reset: function () {
                $(searchSelector).val("");
                $(searchClearSelector).addClass("d-none");
                columnFilterStates = {};
                $('.ag-col-menu-trigger').removeClass('has-filter');
                currentOrderCol = defaultOrderCol;
                currentOrderDir = defaultOrderDir;
                $("#ag_page_size").val(initialPageSize);
                this.reload(1, initialPageSize);
            },

            destroy: function () {
                closeCommonColumnPopup();
                clearTimeout(searchDebounceTimer);

                $(searchSelector).off(".agSearch");
                $(searchClearSelector).off(".agSearchClear");
                $("#ag_page_size").off(".agPageSize");
                $(document).off("click.agPageNav");
                $(exportSelector).off(".agExport");
                $(refreshSelector).off(".agRefresh");

                if (gridApi) {
                    gridApi.destroy();
                    gridApi = null;
                }
            }
        };

        // Wire Event Listeners
        $(searchSelector).off("input.agSearch").on("input.agSearch", function () {
            let val = $(this).val();
            if (val && val.trim().length > 0) {
                $(searchClearSelector).removeClass("d-none");
            } else {
                $(searchClearSelector).addClass("d-none");
            }

            clearTimeout(searchDebounceTimer);
            searchDebounceTimer = setTimeout(function () {
                controller.reload(1, currentPageSize);
            }, 350);
        });

        $(searchClearSelector).off("click.agSearchClear").on("click.agSearchClear", function () {
            $(searchSelector).val("");
            $(this).addClass("d-none");
            controller.reload(1, currentPageSize);
        });

        $("#ag_page_size").off("change.agPageSize").on("change.agPageSize", function () {
            currentPageSize = parseInt($(this).val()) || 10;
            controller.reload(1, currentPageSize);
        });

        $(document).off("click.agPageNav", "#ag_pagination_list .page-link").on("click.agPageNav", "#ag_pagination_list .page-link", function (e) {
            e.preventDefault();
            let parent = $(this).parent();
            if (parent.hasClass("disabled") || parent.hasClass("active")) return;

            let targetPage = parseInt($(this).attr("data-page"));
            let totalPages = Math.ceil(totalRecords / currentPageSize) || 1;

            if (targetPage >= 1 && targetPage <= totalPages) {
                controller.reload(targetPage, currentPageSize);
            }
        });

        $(exportSelector).off("click.agExport").on("click.agExport", function () {
            if (gridApi) {
                gridApi.exportDataAsCsv({
                    fileName: `${exportFileName}_${new Date().toISOString().slice(0, 10)}.csv`
                });
            }
        });

        $(refreshSelector).off("click.agRefresh").on("click.agRefresh", function () {
            let icon = $(this).find("i");
            icon.addClass("ti-spin");
            controller.reset();
            setTimeout(function () {
                icon.removeClass("ti-spin");
            }, 400);
        });

        // Initial Data Load
        controller.reload(1, currentPageSize);

        return controller;
    }

    // Expose globally
    window.initAgDataGrid = initAgDataGrid;

})(window, jQuery);
