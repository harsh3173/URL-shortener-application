import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { Copy, Check, Loader2, Link as LinkIcon } from 'lucide-react';
import { CreateURLRequest, URL } from '@/types';
import { urlApi } from '@/services/api';
import toast from 'react-hot-toast';

interface URLShortenerFormProps {
  onURLCreated?: (url: URL) => void;
}

const URLShortenerForm: React.FC<URLShortenerFormProps> = ({ onURLCreated }) => {
  const [loading, setLoading] = useState(false);
  const [shortenedURL, setShortenedURL] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CreateURLRequest>();

  const onSubmit = async (data: CreateURLRequest) => {
    setLoading(true);
    try {
      const response = await urlApi.createURL(data);
      const shortURL = `${window.location.origin}/${response.data!.short_code}`;
      setShortenedURL(shortURL);
      onURLCreated?.(response.data!);
      toast.success('URL shortened successfully!');
      reset();
    } catch (error: any) {
      const message = error.response?.data?.message || 'Failed to shorten URL';
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = async () => {
    if (!shortenedURL) return;
    
    try {
      await navigator.clipboard.writeText(shortenedURL);
      setCopied(true);
      toast.success('URL copied to clipboard!');
      setTimeout(() => setCopied(false), 2000);
    } catch (error) {
      toast.error('Failed to copy URL');
    }
  };

  return (
    <div className="w-full max-w-4xl mx-auto">
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        <div className="card">
          <div className="card-header">
            <h2 className="text-xl font-semibold text-gray-900 flex items-center">
              <LinkIcon className="mr-2 h-5 w-5" />
              Shorten Your URL
            </h2>
            <p className="text-sm text-gray-600">
              Create a short, memorable link for your URL
            </p>
          </div>
          
          <div className="card-content space-y-4">
            <div>
              <label htmlFor="original_url" className="block text-sm font-medium text-gray-700 mb-2">
                Original URL *
              </label>
              <input
                type="url"
                id="original_url"
                {...register('original_url', {
                  required: 'URL is required',
                  pattern: {
                    value: /^https?:\/\/.+/,
                    message: 'Please enter a valid URL starting with http:// or https://'
                  }
                })}
                className="input w-full"
                placeholder="https://example.com/your-long-url"
              />
              {errors.original_url && (
                <p className="mt-1 text-sm text-red-600">{errors.original_url.message}</p>
              )}
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label htmlFor="custom_alias" className="block text-sm font-medium text-gray-700 mb-2">
                  Custom Alias (Optional)
                </label>
                <input
                  type="text"
                  id="custom_alias"
                  {...register('custom_alias', {
                    pattern: {
                      value: /^[a-zA-Z0-9_-]+$/,
                      message: 'Only letters, numbers, hyphens, and underscores allowed'
                    },
                    minLength: {
                      value: 3,
                      message: 'Alias must be at least 3 characters'
                    },
                    maxLength: {
                      value: 50,
                      message: 'Alias must be less than 50 characters'
                    }
                  })}
                  className="input w-full"
                  placeholder="my-custom-link"
                />
                {errors.custom_alias && (
                  <p className="mt-1 text-sm text-red-600">{errors.custom_alias.message}</p>
                )}
              </div>

              <div>
                <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-2">
                  Title (Optional)
                </label>
                <input
                  type="text"
                  id="title"
                  {...register('title', {
                    maxLength: {
                      value: 200,
                      message: 'Title must be less than 200 characters'
                    }
                  })}
                  className="input w-full"
                  placeholder="My awesome link"
                />
                {errors.title && (
                  <p className="mt-1 text-sm text-red-600">{errors.title.message}</p>
                )}
              </div>
            </div>

            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-2">
                Description (Optional)
              </label>
              <textarea
                id="description"
                {...register('description', {
                  maxLength: {
                    value: 500,
                    message: 'Description must be less than 500 characters'
                  }
                })}
                rows={3}
                className="input w-full resize-none"
                placeholder="Brief description of your link"
              />
              {errors.description && (
                <p className="mt-1 text-sm text-red-600">{errors.description.message}</p>
              )}
            </div>

            <button
              type="submit"
              disabled={loading}
              className="btn btn-primary w-full md:w-auto"
            >
              {loading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Shortening...
                </>
              ) : (
                'Shorten URL'
              )}
            </button>
          </div>
        </div>
      </form>

      {shortenedURL && (
        <div className="card mt-6 animate-slide-up">
          <div className="card-content">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              Your shortened URL is ready!
            </h3>
            <div className="flex items-center space-x-2 p-3 bg-gray-50 rounded-lg">
              <input
                type="text"
                value={shortenedURL}
                readOnly
                className="flex-1 bg-transparent border-none outline-none text-primary-600 font-medium"
              />
              <button
                onClick={copyToClipboard}
                className="btn btn-secondary btn-sm"
              >
                {copied ? (
                  <>
                    <Check className="h-4 w-4 mr-1" />
                    Copied!
                  </>
                ) : (
                  <>
                    <Copy className="h-4 w-4 mr-1" />
                    Copy
                  </>
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default URLShortenerForm;