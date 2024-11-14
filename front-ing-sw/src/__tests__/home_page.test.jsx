import React from 'react';
import { render, screen, fireEvent, within, waitFor, act} from '@testing-library/react';
import { MemoryRouter, useNavigate } from 'react-router-dom';
import HomePage from '../pages/home_page';

// Mock `useNavigate`
const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

describe('HomePage', () => {
  beforeEach(() => {

    jest.clearAllMocks();

    jest.spyOn(global, 'fetch').mockResolvedValue({
      json: jest.fn().mockResolvedValue([
        { id: 1, name: 'John Doe', last_name: 'Doe', username: 'johndoe', email: 'johndoe@example.com', phone: '1234567890', address: '123 Main St' },
        { id: 2, name: 'Jane Smith', last_name: 'Smith', username: 'janesmith', email: 'janesmith@example.com', phone: '0987654321', address: '456 Oak Rd' },
      ]),
      ok: true,
    });
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  test('renders the HomePage component', async () => {
    render(
      <MemoryRouter>
        <HomePage />
      </MemoryRouter>
    );

    // Wait for the users to be fetched
    await screen.findByText('John Doe Doe');
    await screen.findByText('Jane Smith Smith');

    // Check that the search input is rendered
    expect(screen.getByPlaceholderText('Busca por nombre...')).toBeInTheDocument();
  });

  test('filters users based on the search term', async () => {
    render(
      <MemoryRouter>
        <HomePage />
      </MemoryRouter>
    );

    // Wait for the users to be fetched
    await screen.findByText('John Doe Doe');
    await screen.findByText('Jane Smith Smith');

    // Type a search term and verify the filtered users are displayed
    fireEvent.change(screen.getByPlaceholderText('Busca por nombre...'), { target: { value: 'John' } });
    expect(screen.getByText('John Doe Doe')).toBeInTheDocument();
    expect(screen.queryByText('Jane Smith Smith')).not.toBeInTheDocument();
  });

  test('handles user modification', async () => {
    render(
      <MemoryRouter>
        <HomePage />
      </MemoryRouter>
    );

    // Wait for the users to be fetched
    await screen.findByText('John Doe Doe');
    await waitFor(() => screen.getAllByText('Modificar'));

    // Get all "Modificar" buttons and click the first one
    const buttons = screen.getAllByText('Modificar');
    fireEvent.click(buttons[0]);

    // Check if `mockNavigate` was called with the correct argument
    expect(mockNavigate).toHaveBeenCalledWith('/modify-user/1');
  });

  test('handles user elimination', async () => {
    global.window.confirm = jest.fn(() => true);

    // Wrap rendering and any actions that may update state inside `act()`
    await act(async () => {
      render(
        <MemoryRouter>
          <HomePage />
        </MemoryRouter>
      );
    });

    // Wait for the users to be fetched
    await screen.findByText('John Doe Doe');

    // Click the "Eliminar" button and verify the fetch function is called

    const buttons = screen.getAllByText('Eliminar');
    await act(async () => {
      fireEvent.click(buttons[0]);
    });

    expect(global.fetch).toHaveBeenCalledWith(`${process.env.REACT_APP_API_BASE_URL}/user/1`, {
      method: 'DELETE',
    });
  });
});