import React, { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { Loader2 } from 'lucide-react';
import toast from 'react-hot-toast';

const OAuthCallbackPage: React.FC = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { checkAuth } = useAuth();

  useEffect(() => {
    const handleCallback = async () => {
      const code = searchParams.get('code');
      const state = searchParams.get('state');
      const error = searchParams.get('error');

      if (error) {
        toast.error('OAuth authentication failed');
        navigate('/login');
        return;
      }

      if (!code || !state) {
        toast.error('Invalid OAuth callback');
        navigate('/login');
        return;
      }

      try {
        // Wait a moment for the OAuth callback to be processed by backend
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        // Check if the user is now authenticated
        await checkAuth();
        toast.success('Successfully signed in!');
        navigate('/dashboard');
      } catch (error) {
        console.error('OAuth callback error:', error);
        toast.error('Authentication failed');
        navigate('/login');
      }
    };

    handleCallback();
  }, [searchParams, navigate, checkAuth]);

  return (
    <div className="min-h-[80vh] flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8 text-center">
        <div>
          <Loader2 className="mx-auto h-12 w-12 animate-spin text-primary-600" />
          <h2 className="mt-6 text-3xl font-extrabold text-gray-900">
            Completing sign in...
          </h2>
          <p className="mt-2 text-sm text-gray-600">
            Please wait while we finish setting up your account.
          </p>
        </div>
      </div>
    </div>
  );
};

export default OAuthCallbackPage;