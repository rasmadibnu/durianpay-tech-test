import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { beforeEach, describe, expect, it, vi } from "vitest";
import UsersPage from "./UsersPage";

const {
  mockGetUsers,
  mockCreateUser,
  mockUpdateUser,
  mockUpdateUserPassword,
  mockDeleteUser,
  mockToastSuccess,
  mockToastError,
} = vi.hoisted(() => ({
  mockGetUsers: vi.fn(),
  mockCreateUser: vi.fn(),
  mockUpdateUser: vi.fn(),
  mockUpdateUserPassword: vi.fn(),
  mockDeleteUser: vi.fn(),
  mockToastSuccess: vi.fn(),
  mockToastError: vi.fn(),
}));

vi.mock("@/api/endpoints", () => ({
  getUsers: (...args: unknown[]) => mockGetUsers(...args),
  createUser: (...args: unknown[]) => mockCreateUser(...args),
  updateUser: (...args: unknown[]) => mockUpdateUser(...args),
  updateUserPassword: (...args: unknown[]) => mockUpdateUserPassword(...args),
  deleteUser: (...args: unknown[]) => mockDeleteUser(...args),
}));

vi.mock("sonner", () => ({
  toast: {
    success: mockToastSuccess,
    error: mockToastError,
  },
}));

const usersResponse = {
  users: [
    {
      id: "1",
      email: "cs@test.com",
      role: "cs",
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
      <UsersPage />
    </QueryClientProvider>,
  );
}

describe("UsersPage", () => {
  beforeEach(() => {
    mockGetUsers.mockReset();
    mockCreateUser.mockReset();
    mockUpdateUser.mockReset();
    mockUpdateUserPassword.mockReset();
    mockDeleteUser.mockReset();
    mockToastSuccess.mockReset();
    mockToastError.mockReset();

    mockGetUsers.mockResolvedValue(usersResponse);
  });

  it("renders user data", async () => {
    renderPage();

    expect(await screen.findByText("Users")).toBeInTheDocument();
    expect(await screen.findByText("cs@test.com")).toBeInTheDocument();
    expect(screen.getByText("cs")).toBeInTheDocument();
  });

  it("creates a user", async () => {
    const user = userEvent.setup();
    mockCreateUser.mockResolvedValue({});

    renderPage();

    await screen.findByText("Users");
    await user.click(screen.getByRole("button", { name: /add user/i }));
    await user.type(screen.getByPlaceholderText("user@example.com"), "new@test.com");
    await user.type(screen.getByPlaceholderText("Enter password"), "password");
    await user.click(screen.getByRole("button", { name: /^create$/i }));

    await waitFor(() => {
      expect(mockCreateUser).toHaveBeenCalled();
    });
    expect(mockCreateUser.mock.calls[0][0]).toEqual({
      email: "new@test.com",
      password: "password",
      role: "cs",
    });
    expect(mockToastSuccess).toHaveBeenCalledWith("User created successfully");
  });
});
