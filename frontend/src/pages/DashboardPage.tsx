import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { URL } from '@/types';
import { urlApi } from '@/services/api';
import URLShortenerForm from '@/components/URLShortenerForm';
import LoadingSpinner from '@/components/LoadingSpinner';
import { 
  BarChart3, 
  Copy, 
  Trash2, 
  ExternalLink, 
  Calendar,
  MousePointer,
  TrendingUp,
  Link as LinkIcon,
  Plus
} from 'lucide-react';
import { format } from 'date-fns';
import toast from 'react-hot-toast';

const DashboardPage: React.FC = () => {
  const { user } = useAuth();
  const [urls, setUrls] = useState<URL[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [pagination, setPagination] = useState({
    total: 0,
    limit: 10,
    offset: 0,
  });

  useEffect(() => {
    fetchURLs();
  }, [pagination.offset]);

  const fetchURLs = async () => {
    setLoading(true);
    try {
      const response = await urlApi.getUserURLs(pagination.limit, pagination.offset);
      setUrls(response.data!.urls);
      setPagination(prev => ({
        ...prev,
        total: response.data!.total,
      }));
    } catch (error) {
      console.error('Error fetching URLs:', error);
      toast.error('Failed to fetch URLs');
    } finally {
      setLoading(false);
    }
  };

  const handleURLCreated = (newURL: URL) => {
    setUrls(prev => [newURL, ...prev]);
    setPagination(prev => ({ ...prev, total: prev.total + 1 }));
    setShowForm(false);
  };

  const handleCopyURL = async (shortCode: string) => {
    const shortURL = `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'}/${shortCode}`;
    try {
      await navigator.clipboard.writeText(shortURL);
      toast.success('URL copied to clipboard!');
    } catch (error) {
      toast.error('Failed to copy URL');
    }
  };

  const handleDeleteURL = async (id: number) => {
    if (!window.confirm('Are you sure you want to delete this URL?')) return;
    
    try {
      await urlApi.deleteURL(id);
      setUrls(prev => prev.filter(url => url.id !== id));
      setPagination(prev => ({ ...prev, total: prev.total - 1 }));
      toast.success('URL deleted successfully');
    } catch (error) {
      toast.error('Failed to delete URL');
    }
  };

  const totalPages = Math.ceil(pagination.total / pagination.limit);

  if (loading && urls.length === 0) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <div className="flex justify-between items-center mb-6">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
            <p className="text-gray-600">Welcome back, {user?.name}!</p>
          </div>
          <button
            onClick={() => setShowForm(!showForm)}
            className="btn btn-primary"
          >
            <Plus className="h-4 w-4 mr-2" />
            New URL
          </button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <div className="card">
            <div className="card-content">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Total URLs</p>
                  <p className="text-2xl font-bold text-gray-900">{pagination.total}</p>
                </div>
                <LinkIcon className="h-8 w-8 text-primary-600" />
              </div>
            </div>
          </div>
          
          <div className="card">
            <div className="card-content">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Total Clicks</p>
                  <p className="text-2xl font-bold text-gray-900">--</p>
                </div>
                <MousePointer className="h-8 w-8 text-green-600" />
              </div>
            </div>
          </div>
          
          <div className="card">
            <div className="card-content">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Active URLs</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {urls.filter(url => url.is_active).length}
                  </p>
                </div>
                <TrendingUp className="h-8 w-8 text-blue-600" />
              </div>
            </div>
          </div>
          
          <div className="card">
            <div className="card-content">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">This Month</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {urls.filter(url => {
                      const created = new Date(url.created_at);
                      const now = new Date();
                      return created.getMonth() === now.getMonth() && 
                             created.getFullYear() === now.getFullYear();
                    }).length}
                  </p>
                </div>
                <Calendar className="h-8 w-8 text-purple-600" />
              </div>
            </div>
          </div>
        </div>

        {showForm && (
          <div className="mb-8">
            <URLShortenerForm onURLCreated={handleURLCreated} />
          </div>
        )}

        <div className="card">
          <div className="card-header">
            <h2 className="text-xl font-semibold text-gray-900">Your URLs</h2>
          </div>
          <div className="card-content">
            {urls.length === 0 ? (
              <div className="text-center py-12">
                <LinkIcon className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <h3 className="text-lg font-medium text-gray-900 mb-2">No URLs yet</h3>
                <p className="text-gray-600 mb-4">
                  Create your first shortened URL to get started
                </p>
                <button
                  onClick={() => setShowForm(true)}
                  className="btn btn-primary"
                >
                  Create URL
                </button>
              </div>
            ) : (
              <div className="space-y-4">
                {urls.map((url) => (
                  <div key={url.id} className="border rounded-lg p-4 hover:bg-gray-50 transition-colors">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="flex items-center space-x-2 mb-2">
                          <h3 className="text-lg font-medium text-gray-900 truncate">
                            {url.title || 'Untitled'}
                          </h3>
                          {!url.is_active && (
                            <span className="px-2 py-1 text-xs bg-red-100 text-red-800 rounded">
                              Inactive
                            </span>
                          )}
                        </div>
                        
                        <div className="space-y-1 text-sm text-gray-600">
                          <div className="flex items-center space-x-2">
                            <span className="font-medium">Short URL:</span>
                            <code className="bg-gray-100 px-2 py-1 rounded text-primary-600">
                              {import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'}/{url.short_code}
                            </code>
                          </div>
                          
                          <div className="flex items-center space-x-2">
                            <span className="font-medium">Original:</span>
                            <a
                              href={url.original_url}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="text-blue-600 hover:text-blue-800 truncate max-w-md"
                            >
                              {url.original_url}
                            </a>
                          </div>
                          
                          <div className="flex items-center space-x-4 text-xs text-gray-500">
                            <span>Created: {format(new Date(url.created_at), 'MMM d, yyyy')}</span>
                            {url.expires_at && (
                              <span>Expires: {format(new Date(url.expires_at), 'MMM d, yyyy')}</span>
                            )}
                          </div>
                        </div>
                      </div>
                      
                      <div className="flex items-center space-x-2 ml-4">
                        <button
                          onClick={() => handleCopyURL(url.short_code)}
                          className="p-2 text-gray-600 hover:text-primary-600 transition-colors"
                          title="Copy URL"
                        >
                          <Copy className="h-4 w-4" />
                        </button>
                        
                        <Link
                          to={`/analytics/${url.id}`}
                          className="p-2 text-gray-600 hover:text-blue-600 transition-colors"
                          title="View Analytics"
                        >
                          <BarChart3 className="h-4 w-4" />
                        </Link>
                        
                        <a
                          href={url.original_url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="p-2 text-gray-600 hover:text-green-600 transition-colors"
                          title="Visit Original URL"
                        >
                          <ExternalLink className="h-4 w-4" />
                        </a>
                        
                        <button
                          onClick={() => handleDeleteURL(url.id)}
                          className="p-2 text-gray-600 hover:text-red-600 transition-colors"
                          title="Delete URL"
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {totalPages > 1 && (
          <div className="mt-6 flex justify-between items-center">
            <div className="text-sm text-gray-600">
              Showing {pagination.offset + 1} to {Math.min(pagination.offset + pagination.limit, pagination.total)} of {pagination.total} URLs
            </div>
            <div className="flex space-x-2">
              <button
                onClick={() => setPagination(prev => ({ ...prev, offset: prev.offset - prev.limit }))}
                disabled={pagination.offset === 0}
                className="btn btn-outline btn-sm disabled:opacity-50"
              >
                Previous
              </button>
              <button
                onClick={() => setPagination(prev => ({ ...prev, offset: prev.offset + prev.limit }))}
                disabled={pagination.offset + pagination.limit >= pagination.total}
                className="btn btn-outline btn-sm disabled:opacity-50"
              >
                Next
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default DashboardPage;