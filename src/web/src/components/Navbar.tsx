import React from 'react';
import { Link } from '@tanstack/react-router';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { getSession, logOut } from '../lib/api';

const Navbar: React.FC = () => {
  const queryClient = useQueryClient();
  const session = useQuery({ queryKey: ['session'], queryFn: getSession });
  const logout = useMutation({
    mutationFn: logOut,
    onSuccess: () => queryClient.setQueryData(['session'], null),
  });
  return (
    <nav className="fixed top-0 left-0 right-0 z-50 bg-white/80 dark:bg-gray-900/80 backdrop-blur-md border-b border-gray-200 dark:border-gray-800">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-20 items-center">
          {/* Logo and Brand */}
          <div className="flex items-center">
            <Link to="/" className="flex items-center gap-3 group">
              <img
                className="h-10 w-auto transition-transform group-hover:scale-110"
                src="/qmessentials-logo.svg"
                alt="QMEssentials Logo"
              />
              <span className="text-3xl font-extrabold tracking-tight text-gray-900 dark:text-white">
                QMEssentials
              </span>
            </Link>
          </div>

          <div className="flex items-center gap-4">
            {session.data ? (
              <>
                <span className="text-gray-900 dark:text-white">{session.data.userId}</span>
                <button type="button" onClick={() => logout.mutate()} disabled={logout.isPending}
                  className="text-gray-900 dark:text-white">Sign out</button>
              </>
            ) : <Link to="/login" className="text-gray-900 dark:text-white">Sign in</Link>}
          </div>

        </div>
      </div>
    </nav>
  );
};

export default Navbar;
