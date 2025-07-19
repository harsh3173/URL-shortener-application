import React from 'react';
import { Link } from 'react-router-dom';

const Footer: React.FC = () => {
  return (
    <footer className="bg-white border-t">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex justify-center space-x-8">
          <Link to="/" className="text-gray-600 hover:text-primary-600 transition-colors">
            Home
          </Link>
          <Link to="/dashboard" className="text-gray-600 hover:text-primary-600 transition-colors">
            Dashboard
          </Link>
        </div>
        
        <div className="mt-8 pt-8 border-t border-gray-200">
          <p className="text-center text-gray-500 text-sm">
            © 2024 URLShortener. Built with React, TypeScript, Go, and PostgreSQL.
          </p>
        </div>
      </div>
    </footer>
  );
};

export default Footer;