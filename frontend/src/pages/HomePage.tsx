import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import URLShortenerForm from '@/components/URLShortenerForm';
import { ArrowRight } from 'lucide-react';

const HomePage: React.FC = () => {
  const { user } = useAuth();

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="text-center mb-12">
          <h1 className="text-4xl md:text-6xl font-bold text-gray-900 mb-6">
            Shorten URLs with
            <span className="text-primary-600"> Style</span>
          </h1>
          <p className="text-xl text-gray-600 max-w-3xl mx-auto mb-8">
            Create short, memorable links for easy sharing.
          </p>
          {!user && (
            <div className="flex justify-center space-x-4">
              <Link
                to="/register"
                className="btn btn-primary btn-lg"
              >
                Get Started Free
                <ArrowRight className="ml-2 h-5 w-5" />
              </Link>
              <Link
                to="/login"
                className="btn btn-outline btn-lg"
              >
                Sign In
              </Link>
            </div>
          )}
        </div>

        <div className="mb-16">
          <URLShortenerForm />
        </div>

        {!user && (
          <div className="text-center mt-16">
            <Link
              to="/login"
              className="btn btn-primary btn-lg"
            >
              Get Started
              <ArrowRight className="ml-2 h-5 w-5" />
            </Link>
          </div>
        )}
      </div>
    </div>
  );
};

export default HomePage;