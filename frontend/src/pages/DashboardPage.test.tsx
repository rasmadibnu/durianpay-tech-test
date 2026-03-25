import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { beforeEach, describe, expect, it, vi } from "vitest";
import DashboardPage from "./DashboardPage";

const {
  mockUseAuth,
  mockGetPayments,
  mockGetMerchants,
  mockCreatePayment,
  mockUpdatePayment,
  mockReviewPayment,
  mockDeletePayment,
  mockToastSuccess,
  mockToastError,
} = vi.hoisted(() => ({
  mockUseAuth: vi.fn(),
  mockGetPayments: vi.fn(),
  mockGetMerchants: vi.fn(),
  mockCreatePayment: vi.fn(),
  mockUpdatePayment: vi.fn(),
  mockReviewPayment: vi.fn(),
  mockDeletePayment: vi.fn(),
  mockToastSuccess: vi.fn(),
  mockToastError: vi.fn(),
}));

vi.mock("@/context/AuthContext", () => ({
  useAuth: () => mockUseAuth(),
}));

vi.mock("@/services/payment.service", () => ({
  getPayments: (...args: unknown[]) => mockGetPayments(...args),
  getMerchants: (...args: unknown[]) => mockGetMerchants(...args),
  createPayment: (...args: unknown[]) => mockCreatePayment(...args),
  updatePayment: (...args: unknown[]) => mockUpdatePayment(...args),
  reviewPayment: (...args: unknown[]) => mockReviewPayment(...args),
  deletePayment: (...args: unknown[]) => mockDeletePayment(...args),
}));

vi.mock("sonner", () => ({
  toast: {
    success: mockToastSuccess,
    error: mockToastError,
  },
}));

const paymentsResponse = {
  payments: [
    {
      id: "PAY-1",
      merchant_id: 1,
      merchant_name: "Tokopedia",
      amount: "10000",
      status: "processing" as const,
      created_at: "2026-03-25T10:00:00Z",
    },
  ],
  page: 1,
  limit: 10,
  total: 1,
  total_pages: 1,
};

const merchantsResponse = {
  merchants: [
    {
      id: 1,
      name: "Tokopedia",
      created_at: "2026-03-25T10:00:00Z",
      updated_at: "2026-03-25T10:00:00Z",
    },
  ],
  page: 1,
  limit: 100,
  total: 1,
  total_pages: 1,
};

function renderPage() {
  const client = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  return render(
    <QueryClientProvider client={client}>
      <DashboardPage />
    </QueryClientProvider>,
  );
}

describe("DashboardPage", () => {
  beforeEach(() => {
    mockUseAuth.mockReset();
    mockGetPayments.mockReset();
    mockGetMerchants.mockReset();
    mockCreatePayment.mockReset();
    mockUpdatePayment.mockReset();
    mockReviewPayment.mockReset();
    mockDeletePayment.mockReset();
    mockToastSuccess.mockReset();
    mockToastError.mockReset();

    mockGetPayments.mockImplementation(
      (params?: { status?: string; limit?: number }) => {
        if (params?.limit === 1) {
          if (params.status === "completed") {
            return Promise.resolve({
              ...paymentsResponse,
              payments: [],
              total: 4,
            });
          }
          if (params.status === "processing") {
            return Promise.resolve({
              ...paymentsResponse,
              payments: [],
              total: 7,
            });
          }
          if (params.status === "failed") {
            return Promise.resolve({
              ...paymentsResponse,
              payments: [],
              total: 2,
            });
          }
          return Promise.resolve({
            ...paymentsResponse,
            payments: [],
            total: 13,
          });
        }

        if (params?.status === "completed") {
          return Promise.resolve({
            ...paymentsResponse,
            payments: [
              {
                ...paymentsResponse.payments[0],
                id: "PAY-2",
                status: "completed" as const,
              },
            ],
          });
        }

        return Promise.resolve(paymentsResponse);
      },
    );
    mockGetMerchants.mockResolvedValue(merchantsResponse);
  });

  it("renders payment widgets and table data for non-operation users", async () => {
    mockUseAuth.mockReturnValue({ role: "cs" });

    renderPage();

    expect(await screen.findByText("Transactions")).toBeInTheDocument();
    expect(await screen.findByText("Tokopedia")).toBeInTheDocument();
    expect(screen.getByText("Total Payments")).toBeInTheDocument();
    expect(
      screen.queryByRole("button", { name: /add payment/i }),
    ).not.toBeInTheDocument();
  });

  it("creates a payment from the operation flow", async () => {
    const user = userEvent.setup();
    mockUseAuth.mockReturnValue({ role: "operation" });
    mockCreatePayment.mockResolvedValue({});

    renderPage();

    await screen.findByText("Transactions");
    await user.click(screen.getByRole("button", { name: /add payment/i }));
    await user.type(screen.getByPlaceholderText("0.00"), "25000");
    await user.click(screen.getByRole("button", { name: /^create$/i }));

    await waitFor(() => {
      expect(mockCreatePayment).toHaveBeenCalledWith({
        merchant_id: 1,
        amount: "25000",
        status: "processing",
      });
    });
    expect(mockToastSuccess).toHaveBeenCalledWith(
      "Payment created successfully",
    );
  });
});
