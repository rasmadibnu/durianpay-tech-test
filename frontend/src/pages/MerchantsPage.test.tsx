import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { beforeEach, describe, expect, it, vi } from "vitest";
import MerchantsPage from "./MerchantsPage";

const {
  mockGetMerchants,
  mockCreateMerchant,
  mockUpdateMerchant,
  mockDeleteMerchant,
  mockToastSuccess,
  mockToastError,
} = vi.hoisted(() => ({
  mockGetMerchants: vi.fn(),
  mockCreateMerchant: vi.fn(),
  mockUpdateMerchant: vi.fn(),
  mockDeleteMerchant: vi.fn(),
  mockToastSuccess: vi.fn(),
  mockToastError: vi.fn(),
}));

vi.mock("@/services/merchant.service", () => ({
  getMerchants: (...args: unknown[]) => mockGetMerchants(...args),
  createMerchant: (...args: unknown[]) => mockCreateMerchant(...args),
  updateMerchant: (...args: unknown[]) => mockUpdateMerchant(...args),
  deleteMerchant: (...args: unknown[]) => mockDeleteMerchant(...args),
}));

vi.mock("sonner", () => ({
  toast: {
    success: mockToastSuccess,
    error: mockToastError,
  },
}));

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
  limit: 10,
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
      <MerchantsPage />
    </QueryClientProvider>,
  );
}

describe("MerchantsPage", () => {
  beforeEach(() => {
    mockGetMerchants.mockReset();
    mockCreateMerchant.mockReset();
    mockUpdateMerchant.mockReset();
    mockDeleteMerchant.mockReset();
    mockToastSuccess.mockReset();
    mockToastError.mockReset();

    mockGetMerchants.mockResolvedValue(merchantsResponse);
  });

  it("renders merchant data", async () => {
    renderPage();

    expect(await screen.findByText("Merchants")).toBeInTheDocument();
    expect(await screen.findByText("Tokopedia")).toBeInTheDocument();
  });

  it("creates a merchant", async () => {
    const user = userEvent.setup();
    mockCreateMerchant.mockResolvedValue({});

    renderPage();

    await screen.findByText("Merchants");
    await user.click(screen.getByRole("button", { name: /add merchant/i }));
    await user.type(screen.getByPlaceholderText("Merchant name"), "Bukalapak");
    await user.click(screen.getByRole("button", { name: /^create$/i }));

    await waitFor(() => {
      expect(mockCreateMerchant).toHaveBeenCalled();
    });
    expect(mockCreateMerchant.mock.calls[0][0]).toEqual({ name: "Bukalapak" });
    expect(mockToastSuccess).toHaveBeenCalledWith(
      "Merchant created successfully",
    );
  });
});
